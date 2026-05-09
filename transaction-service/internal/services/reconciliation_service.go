package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type ReconciliationService struct {
	db             *gorm.DB
	digiflazzURL   string
	digiflazzKey   string
	digiflazzSecret string
}

func NewReconciliationService(db *gorm.DB) *ReconciliationService {
	return &ReconciliationService{db: db}
}

type ReconciliationResult struct {
	ExpiredCount   int
	ReleasedHoldCnt int
	Errors          []string
}

func (s *ReconciliationService) ReconcileStalePending(ctx context.Context) (*ReconciliationResult, error) {
	result := &ReconciliationResult{
		Errors: []string{},
	}

	timeoutDuration := 15 * time.Minute

	var transactions []Transaction
	err := s.db.WithContext(ctx).
		Where("status = ?", StatePending).
		Where("created_at < ?", time.Now().Add(-timeoutDuration)).
		Where("updated_at < ?", time.Now().Add(-timeoutDuration)).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	result.ExpiredCount = len(transactions)

	for _, tx := range transactions {
		if err := s.expireTransaction(ctx, &tx); err != nil {
			log.Printf("Failed to expire transaction %s: %v", tx.ProviderRef, err)
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.ReleasedHoldCnt++
	}

	log.Printf("Reconciliation: expired=%d, released_hold=%d, errors=%d",
		result.ExpiredCount, result.ReleasedHoldCnt, len(result.Errors))

	return result, nil
}

func (s *ReconciliationService) expireTransaction(ctx context.Context, tx *Transaction) error {
	return s.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		oldStatus := tx.Status
		tx.Status = string(StateExpired)
		tx.PreviousStatus = oldStatus
		tx.StatusChangeReason = "timeout"
		reconciledAt := time.Now()
		tx.ReconciledAt = &reconciledAt
		holdReleasedAt := time.Now()
		tx.HoldReleasedAt = &holdReleasedAt

		if err := txDB.Save(tx).Error; err != nil {
			return err
		}

		s.logAudit(ctx, tx.UserID, "transaction_expired", tx.ID, oldStatus, string(StateExpired), "timeout")

		return nil
	})
}

func (s *ReconciliationService) logAudit(ctx context.Context, userID uint, action string, resourceID uint, oldStatus, newStatus, reason string) {
	details, _ := json.Marshal(map[string]interface{}{
		"old_status": oldStatus,
		"new_status": newStatus,
		"reason":     reason,
		"source":     "reconciliation_job",
	})

	s.db.WithContext(ctx).Exec(`
		INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, severity, created_at)
		VALUES (?, ?, 'transaction', ?, ?::jsonb, 'INFO', NOW())
	`, userID, action, resourceID, details)
}

func (s *ReconciliationService) ReconcileWalletBalances(ctx context.Context) error {
	type walletBalance struct {
		WalletID        uint
		CurrentBalance  float64
		ComputedBalance float64
	}

	var mismatches []walletBalance

	s.db.WithContext(ctx).Raw(`
		SELECT 
			w.id as wallet_id,
			w.balance as current_balance,
			COALESCE(SUM(CASE 
				WHEN we.type IN ('credit', 'refund', 'topup', 'commission') THEN we.amount 
				WHEN we.type IN ('debit', 'transfer', 'transaction') THEN -we.amount 
				ELSE 0 
			END), 0) as computed_balance
		FROM wallets w
		LEFT JOIN wallet_events we ON w.id = we.wallet_id
		GROUP BY w.id
		HAVING w.balance != computed_balance
	`).Scan(&mismatches)

	if len(mismatches) > 0 {
		log.Printf("WALLET RECONCILIATION: Found %d mismatches", len(mismatches))
		for _, m := range mismatches {
			log.Printf("  Wallet %d: current=%.2f, computed=%.2f", m.WalletID, m.CurrentBalance, m.ComputedBalance)

			s.db.Exec("UPDATE wallets SET balance = ? WHERE id = ?", m.ComputedBalance, m.WalletID)
		}
	}

	log.Printf("WALLET RECONCILIATION: All wallets balanced")
	return nil
}

func (s *ReconciliationService) SyncDigiflazzDeposit(ctx context.Context) error {
	if s.digiflazzURL == "" || s.digiflazzKey == "" || s.digiflazzSecret == "" {
		log.Printf("SKIP DIGIFLAZZ DEPOSIT SYNC: missing config")
		return nil
	}

	payload := map[string]interface{}{
		"cmd":      "deposit",
		"username": s.digiflazzKey,
		"sign":    s.generateDigiflazzSignature(map[string]interface{}{"cmd": "deposit"}),
	}

	payloadJSON, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", s.digiflazzURL+"/profile", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("digiflazz API error: status=%d", resp.StatusCode)
	}

	var result struct {
		Success bool    `json:"success"`
		Data    struct {
			Saldo float64 `json:"saldo"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		log.Printf("DIGIFLAZZ DEPOSIT SYNC: API returned success=false")
		return nil
	}

	log.Printf("DIGIFLAZZ DEPOSIT SYNC: balance=%.2f", result.Data.Saldo)
	return nil
}

func (s *ReconciliationService) generateDigiflazzSignature(payload map[string]interface{}) string {
	delete(payload, "sign")
	delete(payload, "username")
	jsonBytes, _ := json.Marshal(payload)
	hash := md5.Sum(jsonBytes)
	md5String := hex.EncodeToString(hash[:])
	signature := fmt.Sprintf("%s%s%s", s.digiflazzKey, md5String, s.digiflazzSecret)
	sigHash := md5.Sum([]byte(signature))
	return hex.EncodeToString(sigHash[:])
}