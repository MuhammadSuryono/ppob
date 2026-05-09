package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yontech/ppob/transaction-service/internal/dto"
	"gorm.io/gorm"
)

type IntegrityService struct {
	db *gorm.DB
}

func NewIntegrityService(db *gorm.DB) *IntegrityService {
	return &IntegrityService{db: db}
}

type IntegrityCheckResult struct {
	WalletID          uint    `json:"wallet_id"`
	BalanceCheck      bool    `json:"balance_check"`
	HoldAmountCheck   bool    `json:"hold_amount_check"`
	EventCountCheck   bool    `json:"event_count_check"`
	DriftAmount       float64 `json:"drift_amount"`
	IsHealthy         bool    `json:"is_healthy"`
	CheckedAt         int64   `json:"checked_at"`
	Issues            []string `json:"issues"`
}

func (s *IntegrityService) CheckWalletIntegrity(ctx context.Context, walletID uint) (*IntegrityCheckResult, error) {
	result := &IntegrityCheckResult{
		WalletID:  walletID,
		CheckedAt: time.Now().Unix(),
		Issues:    []string{},
	}

	var wallet struct {
		ID          uint
		Balance     float64
		HoldAmount  float64
	}
	if err := s.db.Table("wallets").Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	var eventBalance struct {
		Balance float64
	}
	err := s.db.Model(&struct {
		WalletID    uint
		EventType   string
	}{}).
		Select("COALESCE(SUM(CASE WHEN event_type = 'credit' THEN amount ELSE -amount END), 0) as balance").
		Where("wallet_id = ?", walletID).
		Scan(&eventBalance).Error
	if err != nil {
		return nil, err
	}

	drift := wallet.Balance - eventBalance.Balance
	result.DriftAmount = drift

	if drift != 0 {
		result.BalanceCheck = false
		result.Issues = append(result.Issues, fmt.Sprintf("Balance drift detected: %.2f", drift))
	} else {
		result.BalanceCheck = true
	}

	var negativeHoldCount int64
	s.db.Model(&struct {
		WalletID uint
		Status   string
		Amount   float64
	}{}).
		Where("wallet_id = ? AND status = ? AND amount < 0", walletID, "active").
		Count(&negativeHoldCount)

	if negativeHoldCount > 0 {
		result.HoldAmountCheck = false
		result.Issues = append(result.Issues, "Negative hold amount detected")
	} else {
		result.HoldAmountCheck = true
	}

	var eventCount int64
	s.db.Model(&struct {
		WalletID uint
	}{}).
		Where("wallet_id = ?", walletID).
		Count(&eventCount)

	if eventCount == 0 && wallet.Balance != 0 {
		result.EventCountCheck = false
		result.Issues = append(result.Issues, "Wallet has balance but no events")
	} else {
		result.EventCountCheck = true
	}

	result.IsHealthy = result.BalanceCheck && result.HoldAmountCheck && result.EventCountCheck

	return result, nil
}

func (s *IntegrityService) RunDailyIntegrityCheck(ctx context.Context) ([]IntegrityCheckResult, error) {
	var walletIDs []uint
	s.db.Model(&struct {
		ID uint
	}{}).Pluck("id", &walletIDs)

	results := make([]IntegrityCheckResult, 0)
	for _, walletID := range walletIDs {
		result, err := s.CheckWalletIntegrity(ctx, walletID)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}

	return results, nil
}

func (s *IntegrityService) CheckTransactionIntegrity(transactionID string) (*dto.FinancialIntegrityCheckResponse, error) {
	resp := &dto.FinancialIntegrityCheckResponse{
		CheckedAt: time.Now().Unix(),
	}

	var transaction struct {
		ID        uint
		Status    string
		Amount    float64
		Margin    float64
		Completed bool
	}
	if err := s.db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		return nil, err
	}

	resp.TransactionID = transactionID

	if transaction.Status == "completed" {
		var walletDebit float64
		s.db.Model(&struct {
			ReferenceID   string
			ReferenceType string
			EventType     string
		}{}).
			Where("reference_id = ? AND reference_type = ? AND event_type = ?", transactionID, "transaction", "debit").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&walletDebit)

		if walletDebit > 0 {
			resp.BalanceCheck = true
		} else {
			resp.BalanceCheck = false
		}

		if transaction.Margin >= 0 {
			resp.MarginCheck = true
		} else {
			resp.MarginCheck = false
		}

		resp.IsHealthy = resp.BalanceCheck && resp.MarginCheck
	}

	return resp, nil
}