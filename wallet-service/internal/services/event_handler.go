package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yontech/ppob/wallet-service/internal/models"
)

type WalletEventHandler struct {
	walletService *WalletService
}

func NewWalletEventHandler(walletService *WalletService) *WalletEventHandler {
	return &WalletEventHandler{
		walletService: walletService,
	}
}

func (h *WalletEventHandler) HandleEvent(ctx context.Context, eventType string, payload string) error {
	log.Printf("Handling event: %s", eventType)

	switch eventType {
	case "user.registered":
		return h.handleUserRegistered(ctx, payload)
	case "transaction.success":
		return h.handleTransactionSuccess(ctx, payload)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func (h *WalletEventHandler) handleUserRegistered(ctx context.Context, payload string) error {
	var data struct {
		UserID uint `json:"user_id"`
	}

	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return fmt.Errorf("failed to unmarshal user.registered payload: %w", err)
	}

	// Check if wallet already exists (might be created by DB trigger)
	existing, _ := h.walletService.walletRepo.FindByUserID(data.UserID)
	if existing != nil {
		log.Printf("Wallet already exists for user %d, skipping creation", data.UserID)
		return nil
	}

	// Create initial wallet for user
	wallet := &models.Wallet{
		UserID:  data.UserID,
		Balance: 0,
		Status:  "active",
	}

	return h.walletService.walletRepo.Create(wallet)
}

func (h *WalletEventHandler) handleTransactionSuccess(ctx context.Context, payload string) error {
	var data struct {
		TransactionID   string  `json:"transaction_id"`
		UserID          uint    `json:"user_id"`
		Amount          float64 `json:"amount"`
		ProductCode     string  `json:"product_code"`
		StaffCommission float64 `json:"staff_commission"`
	}

	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return fmt.Errorf("failed to unmarshal transaction.success payload: %w", err)
	}

	if data.StaffCommission > 0 {
		log.Printf("Crediting commission of %.2f to user %d for transaction %s", data.StaffCommission, data.UserID, data.TransactionID)
		// We use level 1 as default for direct staff commission
		return h.walletService.CreateCommission(ctx, data.UserID, data.StaffCommission, data.TransactionID, "MarginShare", 1)
	}

	log.Printf("Transaction %s success event received for user %d (no commission)", data.TransactionID, data.UserID)
	return nil
}
