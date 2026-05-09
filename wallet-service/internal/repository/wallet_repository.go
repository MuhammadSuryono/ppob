package repository

import (
	"time"

	"github.com/yontech/ppob/wallet-service/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) FindByUserIDForUpdate(tx *gorm.DB, userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *WalletRepository) UpdateBalanceWithTx(tx *gorm.DB, walletID uint, amount float64, txType string) error {
	return tx.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *WalletRepository) UpdateHoldAmountWithTx(tx *gorm.DB, walletID uint, amount float64, operation string) error {
	expr := "hold_amount + ?"
	if operation == "sub" {
		expr = "hold_amount - ?"
	}
	return tx.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("hold_amount", gorm.Expr(expr, amount)).Error
}

func (r *WalletRepository) AtomicCredit(walletID uint, amount float64, referenceID, referenceType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := r.FindByIDForUpdate(tx, walletID)
		if err != nil {
			return err
		}

		balanceBefore := wallet.Balance

		event := &models.WalletEvent{
			WalletID:      walletID,
			EventType:     "credit",
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:   balanceBefore + amount,
			ReferenceID:   referenceID,
			ReferenceType: referenceType,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		return tx.Save(wallet).Error
	})
}

func (r *WalletRepository) AtomicDebit(walletID uint, amount float64, referenceID, referenceType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := r.FindByIDForUpdate(tx, walletID)
		if err != nil {
			return err
		}

		available := wallet.Balance - wallet.HoldAmount
		if available < amount {
			return ErrInsufficientFunds
		}

		balanceBefore := wallet.Balance

		event := &models.WalletEvent{
			WalletID:      walletID,
			EventType:     "debit",
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:   balanceBefore - amount,
			ReferenceID:   referenceID,
			ReferenceType: referenceType,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		wallet.Balance -= amount
		return tx.Save(wallet).Error
	})
}

func (r *WalletRepository) AtomicTransfer(fromWalletID, toWalletID uint, amount float64, referenceID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		fromWallet, err := r.FindByIDForUpdate(tx, fromWalletID)
		if err != nil {
			return err
		}

		available := fromWallet.Balance - fromWallet.HoldAmount
		if available < amount {
			return ErrInsufficientFunds
		}

		toWallet, err := r.FindByIDForUpdate(tx, toWalletID)
		if err != nil {
			return err
		}

		debitEvent := &models.WalletEvent{
			WalletID:      fromWalletID,
			EventType:     "debit",
			Amount:        amount,
			BalanceBefore: fromWallet.Balance,
			BalanceAfter:   fromWallet.Balance - amount,
			ReferenceID:   referenceID,
			ReferenceType: "transfer",
		}
		if err := tx.Create(debitEvent).Error; err != nil {
			return err
		}
		fromWallet.Balance -= amount
		if err := tx.Save(fromWallet).Error; err != nil {
			return err
		}

		creditEvent := &models.WalletEvent{
			WalletID:      toWalletID,
			EventType:     "credit",
			Amount:        amount,
			BalanceBefore: toWallet.Balance,
			BalanceAfter:   toWallet.Balance + amount,
			ReferenceID:   referenceID,
			ReferenceType: "transfer",
		}
		if err := tx.Create(creditEvent).Error; err != nil {
			return err
		}
		toWallet.Balance += amount
		return tx.Save(toWallet).Error
	})
}

var ErrInsufficientFunds = &WalletError{Message: "insufficient funds"}

type WalletError struct {
	Message string
}

func (e *WalletError) Error() string {
	return e.Message
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.WalletEvent) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) FindByWalletID(walletID uint, limit, offset int) ([]models.WalletEvent, error) {
	var events []models.WalletEvent
	err := r.db.Where("wallet_id = ?", walletID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&events).Error
	return events, err
}

func (r *EventRepository) GetBalanceByEvents(walletID uint) (float64, error) {
	var result struct {
		Total float64
	}
	err := r.db.Model(&models.WalletEvent{}).
		Where("wallet_id = ?", walletID).
		Select("COALESCE(SUM(CASE WHEN event_type = 'credit' THEN amount ELSE -amount END), 0) as total").
		Scan(&result).Error
	return result.Total, err
}

func (r *EventRepository) GetEventSummary(walletID uint) (map[string]float64, error) {
	var results []struct {
		EventType string
		Total     float64
	}
	err := r.db.Model(&models.WalletEvent{}).
		Where("wallet_id = ?", walletID).
		Group("event_type").
		Select("event_type, SUM(amount) as total").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	summary := make(map[string]float64)
	for _, r := range results {
		summary[r.EventType] = r.Total
	}
	return summary, nil
}

func (r *EventRepository) GetEventsByDateRange(walletID uint, start, end time.Time) ([]models.WalletEvent, error) {
	var events []models.WalletEvent
	err := r.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?", walletID, start, end).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}

type HoldRepository struct {
	db *gorm.DB
}

func NewHoldRepository(db *gorm.DB) *HoldRepository {
	return &HoldRepository{db: db}
}

func (r *HoldRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *HoldRepository) Create(hold *models.Hold) error {
	return r.db.Create(hold).Error
}

func (r *HoldRepository) CreateWithTx(tx *gorm.DB, hold *models.Hold) error {
	return tx.Create(hold).Error
}

func (r *HoldRepository) FindByReference(referenceID, referenceType string) (*models.Hold, error) {
	var hold models.Hold
	err := r.db.Where("reference_id = ? AND reference_type = ? AND status = ?", referenceID, referenceType, "active").First(&hold).Error
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

func (r *HoldRepository) FindByReferenceForUpdate(tx *gorm.DB, referenceID, referenceType string) (*models.Hold, error) {
	var hold models.Hold
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("reference_id = ? AND reference_type = ? AND status = ?", referenceID, referenceType, "active").First(&hold).Error
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

func (r *HoldRepository) Release(hold *models.Hold) error {
	return r.db.Model(hold).Updates(map[string]interface{}{"status": "released", "released_at": time.Now()}).Error
}

func (r *HoldRepository) ReleaseWithTx(tx *gorm.DB, hold *models.Hold) error {
	return tx.Model(hold).Updates(map[string]interface{}{"status": "released", "released_at": time.Now()}).Error
}

func (r *HoldRepository) UpdateWalletHoldAmount(walletID uint, amount float64, operation string) error {
	var wallet models.Wallet
	if err := r.db.First(&wallet, walletID).Error; err != nil {
		return err
	}

	if operation == "add" {
		wallet.HoldAmount += amount
	} else {
		wallet.HoldAmount -= amount
	}

	return r.db.Save(&wallet).Error
}

func (r *HoldRepository) UpdateWalletHoldAmountWithTx(tx *gorm.DB, walletID uint, amount float64, operation string) error {
	expr := "hold_amount + ?"
	if operation == "sub" {
		expr = "hold_amount - ?"
	}
	return tx.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("hold_amount", gorm.Expr(expr, amount)).Error
}

func (r *HoldRepository) AtomicPlaceHold(walletID uint, amount float64, referenceID, referenceType string, expiresAt *time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := NewWalletRepository(r.db).FindByIDForUpdate(tx, walletID)
		if err != nil {
			return err
		}

		available := wallet.Balance - wallet.HoldAmount
		if available < amount {
			return ErrInsufficientFunds
		}

		hold := &models.Hold{
			WalletID:      walletID,
			Amount:        amount,
			ReferenceID:   referenceID,
			ReferenceType: referenceType,
			Status:        "active",
			ExpiresAt:     expiresAt,
		}
		if err := tx.Create(hold).Error; err != nil {
			return err
		}

		return tx.Model(&models.Wallet{}).
			Where("id = ?", walletID).
			Update("hold_amount", gorm.Expr("hold_amount + ?", amount)).Error
	})
}

func (r *HoldRepository) AtomicReleaseHold(referenceID, referenceType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		hold, err := r.FindByReferenceForUpdate(tx, referenceID, referenceType)
		if err != nil {
			return err
		}

		if hold.ExpiresAt != nil && hold.ExpiresAt.Before(time.Now()) {
			return &HoldError{Message: "hold expired"}
		}

		if err := r.ReleaseWithTx(tx, hold); err != nil {
			return err
		}

		return tx.Model(&models.Wallet{}).
			Where("id = ?", hold.WalletID).
			Update("hold_amount", gorm.Expr("hold_amount - ?", hold.Amount)).Error
	})
}

var ErrHoldExpired = &HoldError{Message: "hold expired"}

type HoldError struct {
	Message string
}

func (e *HoldError) Error() string {
	return e.Message
}