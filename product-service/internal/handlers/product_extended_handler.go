package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/product-service/internal/services"
)

type SyncHandler struct {
	syncService *services.ProductSyncService
}

func NewSyncHandler(syncService *services.ProductSyncService) *SyncHandler {
	return &SyncHandler{syncService: syncService}
}

func (h *SyncHandler) SyncPrepaid(c *gin.Context) {
	err := h.syncService.SyncPrepaidProducts(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Prepaid products synced successfully",
		"synced_at":   "now",
	})
}

func (h *SyncHandler) SyncPostpaid(c *gin.Context) {
	err := h.syncService.SyncPostpaidProducts(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Postpaid products synced successfully",
		"synced_at":   "now",
	})
}

func (h *SyncHandler) GetSyncStatus(c *gin.Context) {
	prepaidSync, _ := h.syncService.GetLastSyncTime(c.Request.Context(), "prepaid")
	postpaidSync, _ := h.syncService.GetLastSyncTime(c.Request.Context(), "postpaid")

	c.JSON(http.StatusOK, gin.H{
		"prepaid": gin.H{
			"last_sync": prepaidSync.Unix(),
		},
		"postpaid": gin.H{
			"last_sync": postpaidSync.Unix(),
		},
	})
}

type PriceValidationHandler struct {
	priceValidationService *services.PriceValidationService
}

func NewPriceValidationHandler(priceValidationService *services.PriceValidationService) *PriceValidationHandler {
	return &PriceValidationHandler{priceValidationService: priceValidationService}
}

func (h *PriceValidationHandler) ValidatePrice(c *gin.Context) {
	productCode := c.Query("product_code")
	sellingPriceStr := c.Query("selling_price")

	if productCode == "" || sellingPriceStr == "" {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"fields": "product_code, selling_price"}))
		return
	}

	sellingPrice, err := strconv.ParseFloat(sellingPriceStr, 64)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	result := h.priceValidationService.ValidatePrice(productCode, sellingPrice)

	if !result.Valid {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_PRICE_BELOW_COST", map[string]interface{}{"product_code": result.ProductCode}))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PriceValidationHandler) BatchValidatePrice(c *gin.Context) {
	var pricing map[string]float64
	if err := c.ShouldBindJSON(&pricing); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	results := h.priceValidationService.BatchValidate(pricing)

	c.JSON(http.StatusOK, gin.H{"results": results})
}