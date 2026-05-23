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
		TransactionID string  `json:"transaction_id"`
		UserID        uint    `json:"user_id"`
		Amount        float64 `json:"amount"`
		ProductCode   string  `json:"product_code"`
	}

	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return fmt.Errorf("failed to unmarshal transaction.success payload: %w", err)
	}

	// Logic for commission distribution would go here if not handled by Transaction Service.
	// As per blueprint Section 4.2: "Wallet Service: (Event) Credits staff commission."
	// Note: Currently transaction-service's CommissionWorker also does this via gRPC. 
	// To fully align with Section 4.2, we should move commission calculation/distribution here.
	
	log.Printf("Transaction %s success event received for user %d", data.TransactionID, data.UserID)
	return nil
}
