package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/services"
)

type MarginHandler struct {
	marginService *services.MarginService
}

func NewMarginHandler(marginService *services.MarginService) *MarginHandler {
	return &MarginHandler{marginService: marginService}
}

func (h *MarginHandler) SetStaffMargin(c *gin.Context) {
	mitraID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_TOKEN_EXPIRED", nil))
		return
	}

	var req dto.SetMarginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	setting, err := h.marginService.SetStaffMargin(mitraID.(uint), req.StaffID, req.SchemeType, req.Value)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Margin settings updated",
		"data":    setting,
	})
}

func (h *MarginHandler) GetStaffMargin(c *gin.Context) {
	staffID := c.Param("id")
	var staffIDUint uint
	_, err := dto.ParseUint(staffID, &staffIDUint)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	setting, err := h.marginService.GetStaffMarginSettings(staffIDUint)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (h *MarginHandler) GetStaffProductOverrides(c *gin.Context) {
	staffID := c.Param("id")
	var staffIDUint uint
	_, err := dto.ParseUint(staffID, &staffIDUint)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	overrides, err := h.marginService.GetStaffProductOverrides(staffIDUint)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": overrides})
}

func (h *MarginHandler) SetProductMarginOverride(c *gin.Context) {
	mitraID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_TOKEN_EXPIRED", nil))
		return
	}

	var req dto.SetProductMarginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	override, err := h.marginService.SetProductMarginOverride(mitraID.(uint), req.StaffID, req.ProductID, req.ProductCode, req.Value, req.SchemeType)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product margin override updated",
		"data":    override,
	})
}