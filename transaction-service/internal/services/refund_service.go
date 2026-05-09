package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/yontech/ppob/transaction-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrRefundFailed      = errors.New("refund processing failed")
	ErrRefundAlreadyDone = errors.New("refund already processed")
)

type RefundService struct {
	db              *gorm.DB
	marginService   *MarginService
	walletClient    WalletClientInterface
}

type WalletClientInterface interface {
	CreditWallet(walletID uint, amount float64, referenceID, referenceType string) error
	DebitWallet(walletID uint, amount float64, referenceID, referenceType string) error
}

func NewRefundService(db *gorm.DB, marginService *MarginService, walletClient WalletClientInterface) *RefundService {
	return &RefundService{
		db:            db,
		marginService: marginService,
		walletClient:  walletClient,
	}
}

func (s *RefundService) ProcessRefund(transactionID string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var transaction models.Transaction
	err := tx.Where("transaction_id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return ErrRefundFailed
	}

	if transaction.Status == "Refunded" {
		return ErrRefundAlreadyDone
	}

	walletID := transaction.WalletID
	refundAmount := transaction.SellingPrice

	if walletID != nil {
		err = s.walletClient.CreditWallet(*walletID, refundAmount, transactionID, "refund")
		if err != nil {
			tx.Rollback()
			return ErrRefundFailed
		}
	}

	var commissions []models.Commission
	err = tx.Where("transaction_id = ? AND amount > 0", transactionID).Find(&commissions).Error
	if err != nil {
		tx.Rollback()
		return ErrRefundFailed
	}

	for _, commission := range commissions {
		negativeCommission := &models.Commission{
			StaffID:       commission.StaffID,
			TransactionID: transaction.ID,
			Amount:        -commission.Amount,
			Type:          "refund",
			SchemeUsed:    commission.SchemeUsed,
			MarginAmount:  -commission.MarginAmount,
			PaidAt:        nil,
		}
		err = tx.Create(negativeCommission).Error
		if err != nil {
			tx.Rollback()
			return ErrRefundFailed
		}

		walletID := commission.WalletID
		if walletID != nil && *walletID != 0 {
			err = s.walletClient.DebitWallet(*walletID, commission.Amount, fmt.Sprintf("%d", transaction.ID), "commission_refund")
			if err != nil {
				tx.Rollback()
				return ErrRefundFailed
			}
		}
	}

	transaction.Status = "Refunded"
	transaction.PreviousStatus = "Success"
	transaction.UpdatedAt = time.Now()
	err = tx.Save(&transaction).Error
	if err != nil {
		tx.Rollback()
		return ErrRefundFailed
	}

	return tx.Commit().Error
}

func (s *RefundService) GetRefundsByUser(userID uint, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.db.Where("user_id = ? AND status = ? AND created_at BETWEEN ? AND ?", 
		userID, "Refunded", startDate, endDate).Find(&transactions).Error
	return transactions, err
}