package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/integration-service/internal/dto"
	"github.com/yontech/ppob/integration-service/internal/services"
)

type IntegrationHandler struct {
	integrationService *services.IntegrationService
	digiflazzClient    *services.DigiflazzClient
	compensationService *services.CompensationService
}

func NewIntegrationHandler(integrationService *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{integrationService: integrationService}
}

func NewIntegrationHandlerWithClient(
	integrationService *services.IntegrationService,
	digiflazzClient *services.DigiflazzClient,
	compensationService *services.CompensationService,
) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService:  integrationService,
		digiflazzClient:     digiflazzClient,
		compensationService: compensationService,
	}
}

func (h *IntegrationHandler) InitiateDigiflazzTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.DigiflazzTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.integrationService.InitiateDigiflazzTransaction(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("DIGIFLAZZ_GENERAL_FAILURE", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *IntegrationHandler) HandleDigiflazzCallback(c *gin.Context) {
	var req dto.DigiflazzCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.integrationService.HandleDigiflazzCallback(c.Request.Context(), &req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Callback processed"})
}

func (h *IntegrationHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Digiflazz-Signature")
	timestamp := c.GetHeader("X-Digiflazz-Timestamp")

	bodyBytes, err := c.GetRawData()
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if h.digiflazzClient != nil && signature != "" && timestamp != "" {
		if !h.digiflazzClient.VerifyWebhookSignature(string(bodyBytes), timestamp, signature) {
			errors.RespondWithError(c, errors.NewAppError("DIGIFLAZZ_INVALID_SIGNATURE", nil))
			return
		}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	lockKey := "webhook:processing:" + payload["ref_id"].(string)
	ctx := context.Background()

	locked, err := h.integrationService.TryLockWebhook(ctx, lockKey)
	if err != nil || !locked {
		c.JSON(http.StatusOK, gin.H{"message": "Webhook already processed"})
		return
	}
	defer h.integrationService.ReleaseWebhookLock(ctx, lockKey)

	refID, _ := payload["ref_id"].(string)
	trxID, _ := payload["trx_id"].(string)
	if trxID == "" {
		trxID = refID
	}

	status, _ := payload["status"].(string)
	scCode, _ := payload["sc_code"].(string)
	message, _ := payload["message"].(string)

	webhookStatus, _ := services.MapRCToStatus(scCode)
	if webhookStatus == "" {
		webhookStatus = mapStatus(status)
	}

	updateReq := &dto.UpdateStatusRequest{
		Status:        webhookStatus,
		ProviderRef:  trxID,
		ProviderStatus: status,
		Message:       message,
	}

	if err := h.integrationService.UpdateTransactionStatus(ctx, refID, updateReq); err != nil {
		if h.compensationService != nil {
			h.compensationService.CreateJob(ctx, refID, "webhook_update", map[string]interface{}{
				"status":         webhookStatus,
				"provider_ref":   trxID,
				"provider_status": status,
				"message":        message,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

func mapStatus(status string) string {
	switch status {
	case "Sukses", "Success":
		return "success"
	case "Pending", "Processing":
		return "pending"
	default:
		return "failed"
	}
}

func (h *IntegrationHandler) ListProviders(c *gin.Context) {
	resp, err := h.integrationService.ListProviders(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"providers": resp})
}

func (h *IntegrationHandler) GetErrorCatalog(c *gin.Context) {
	catalog := services.NewErrorCatalog()
	errs := catalog.GetAllErrors()

	c.JSON(http.StatusOK, gin.H{"errors": errs})
}

func (h *IntegrationHandler) GetCompensationJobs(c *gin.Context) {
	if h.compensationService == nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_MAINTENANCE", nil))
		return
	}

	pending, err := h.compensationService.GetPendingJobs(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	scheduled, err := h.compensationService.GetScheduledJobs(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pending":   pending,
		"scheduled": scheduled,
	})
}

func (h *IntegrationHandler) GetDeadLetterQueue(c *gin.Context) {
	if h.compensationService == nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_MAINTENANCE", nil))
		return
	}

	deadLetter, err := h.compensationService.GetDeadLetterJobs(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"dead_letter": deadLetter})
}