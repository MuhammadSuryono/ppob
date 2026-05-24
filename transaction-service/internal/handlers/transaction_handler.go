package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/services"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) InitiateTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	idempotencyKey := c.GetHeader("Idempotency-Key")

	var req dto.InitiateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.transactionService.InitiateTransaction(c.Request.Context(), userID.(uint), &dto.CreateTransactionRequest{
		ProductCode:    req.ProductCode,
		CustomerNumber: req.CustomerNumber,
		Amount:         req.Amount,
		SellingPrice:   req.SellingPrice,
		AuthorizeID:    req.AuthorizeID,
	}, idempotencyKey)

	if err != nil {
		if err == services.ErrIdempotencyKeyUsed {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_ALREADY_PROCESSING", nil))
			return
		}
		if err == services.ErrUnauthorizedTransaction {
			errors.RespondWithError(c, errors.NewAppError("AUTH_AUTHORIZE_INVALID", nil))
			return
		}
		// Debug: return real error message
		errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")

	resp, err := h.transactionService.InitiateTransaction(c.Request.Context(), userID.(uint), &req, idempotencyKey)
	if err != nil {
		if err == services.ErrUnauthorizedTransaction {
			errors.RespondWithError(c, errors.NewAppError("AUTH_AUTHORIZE_INVALID", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.transactionService.GetTransaction(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionID := c.Param("id")

	resp, err := h.transactionService.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.transactionService.UpdateTransactionStatus(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == services.ErrInvalidState {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_CANCELLED", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.transactionService.ListTransactions(c.Request.Context(), userID.(uint), status, startDate, endDate, limit, offset)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.transactionService.GetTransactionHistory(c.Request.Context(), userID.(uint), cursor, limit)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) ProcessWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.transactionService.ProcessWebhook(c.Request.Context(), payload)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("DIGIFLAZZ_UNKNOWN_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.CancelTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = "user_requested"
	}

	userID, _ := c.Get("user_id")
	resp, err := h.transactionService.CancelTransaction(c.Request.Context(), uint(id), userID.(uint), req.Reason)
	if err != nil {
		if err == services.ErrInvalidState {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_CANCELLED", nil))
			return
		}
		if err == services.ErrTransactionNotFound {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetReports returns aggregated KPIs, sales trend, and staff performance
func (h *TransactionHandler) GetReports(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// Only Mitra or Admin can access reports
	if role != "mitra" && role != "admin" {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	var reportUserID uint
	if role == "mitra" {
		reportUserID = userID.(uint)
	} else {
		// Admin can filter by user_id query param
		if uid := c.Query("user_id"); uid != "" {
			if parsed, err := strconv.ParseUint(uid, 10, 32); err == nil {
				reportUserID = uint(parsed)
			}
		}
	}

	resp, err := h.transactionService.GetReports(c.Request.Context(), startDate, endDate, reportUserID)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}
