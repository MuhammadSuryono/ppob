package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/wallet-service/internal/dto"
	"github.com/yontech/ppob/wallet-service/internal/services"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	resp, err := h.walletService.GetBalance(c.Request.Context(), userID.(uint))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *WalletHandler) GetBalanceByEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	resp, err := h.walletService.GetBalanceByEvents(c.Request.Context(), userID.(uint))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *WalletHandler) PlaceHold(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.HoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.walletService.PlaceHold(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == services.ErrInsufficientFund {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_INSUFFICIENT_BALANCE", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *WalletHandler) ReleaseHold(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.ReleaseHoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.ReleaseHold(c.Request.Context(), uint(id), &req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hold released successfully"})
}

func (h *WalletHandler) Debit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.DebitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.Debit(c.Request.Context(), uint(id), &req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Debited successfully"})
}

func (h *WalletHandler) Credit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.CreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.Credit(c.Request.Context(), uint(id), &req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credited successfully"})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.Transfer(c.Request.Context(), userID.(uint), req.ToUserID, req.Amount, req.ReferenceID); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (h *WalletHandler) TopUpStaff(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.TopUpStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.TopUpStaff(c.Request.Context(), userID.(uint), req.StaffUserID, req.Amount, req.ReferenceID); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Top-up successful"})
}

// TopUp allows Mitra to top up their own wallet
func (h *WalletHandler) TopUp(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	var req dto.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.walletService.TopUp(c.Request.Context(), userID.(uint), req.Amount); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Top-up successful"})
}

func (h *WalletHandler) GetEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.walletService.GetEvents(c.Request.Context(), userID.(uint), limit, offset)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": resp})
}

func (h *WalletHandler) Reconcile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.walletService.Reconcile(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

type HoldForTransactionRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type DebitForTransactionRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *WalletHandler) PlaceHoldForTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req HoldForTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "transaction_id"}))
		return
	}

	err := h.walletService.PlaceHoldForTransaction(c.Request.Context(), userID.(uint), req.Amount, transactionID)
	if err != nil {
		if err == services.ErrInsufficientFund {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_INSUFFICIENT_BALANCE", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_HOLD_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hold placed for transaction"})
}

func (h *WalletHandler) ReleaseHoldForTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")
	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "transaction_id"}))
		return
	}

	err := h.walletService.ReleaseHoldForTransaction(c.Request.Context(), userID.(uint), transactionID)
	if err != nil {
		if err == services.ErrHoldNotFound {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_HOLD_NOT_FOUND", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hold released for transaction"})
}

func (h *WalletHandler) DebitForTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req DebitForTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "transaction_id"}))
		return
	}

	err := h.walletService.DebitForTransaction(c.Request.Context(), userID.(uint), req.Amount, transactionID)
	if err != nil {
		if err == services.ErrInsufficientFund {
			errors.RespondWithError(c, errors.NewAppError("TRANSACTION_INSUFFICIENT_BALANCE", nil))
			return
		}
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Debited for transaction"})
}