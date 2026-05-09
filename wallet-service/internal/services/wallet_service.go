package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/wallet-service/config"
	"github.com/yontech/ppob/wallet-service/internal/dto"
	"github.com/yontech/ppob/wallet-service/internal/models"
	"github.com/yontech/ppob/wallet-service/internal/repository"
)

var (
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrInsufficientFund = errors.New("insufficient funds")
	ErrHoldNotFound     = errors.New("hold not found")
	ErrHoldExpired      = errors.New("hold expired")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidUser     = errors.New("invalid user for this operation")
	ErrDailyLimitExceeded = errors.New("daily limit exceeded")
)

type WalletService struct {
	walletRepo    *repository.WalletRepository
	eventRepo     *repository.EventRepository
	holdRepo      *repository.HoldRepository
	redis         *redis.Client
	cfg           *config.Config
}

func NewWalletService(
	walletRepo *repository.WalletRepository,
	eventRepo *repository.EventRepository,
	redis *redis.Client,
	cfg *config.Config,
) *WalletService {
	holdRepo := repository.NewHoldRepository(walletRepo.GetDB())
	return &WalletService{
		walletRepo: walletRepo,
		eventRepo:  eventRepo,
		holdRepo:   holdRepo,
		redis:      redis,
		cfg:        cfg,
	}
}

func (s *WalletService) GetBalance(ctx context.Context, userID uint) (*dto.BalanceResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	return &dto.BalanceResponse{
		WalletID:   wallet.ID,
		UserID:     wallet.UserID,
		Balance:    wallet.Balance,
		HoldAmount: wallet.HoldAmount,
		Available:  wallet.Balance - wallet.HoldAmount,
		Currency:   wallet.Currency,
	}, nil
}

func (s *WalletService) GetBalanceByEvents(ctx context.Context, userID uint) (*dto.BalanceResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	calculatedBalance, err := s.eventRepo.GetBalanceByEvents(wallet.ID)
	if err != nil {
		return nil, err
	}

	return &dto.BalanceResponse{
		WalletID:      wallet.ID,
		UserID:        wallet.UserID,
		Balance:       wallet.Balance,
		HoldAmount:    wallet.HoldAmount,
		Available:     wallet.Balance - wallet.HoldAmount,
		Currency:      wallet.Currency,
		CalculatedBalance: calculatedBalance,
	}, nil
}

func (s *WalletService) PlaceHold(ctx context.Context, userID uint, req *dto.HoldRequest) (*dto.HoldResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	var expiresAt *time.Time
	if req.ExpiresAt > 0 {
		exp := time.Unix(req.ExpiresAt, 0)
		expiresAt = &exp
	}

	if err := s.holdRepo.AtomicPlaceHold(wallet.ID, req.Amount, req.ReferenceID, req.ReferenceType, expiresAt); err != nil {
		if err == repository.ErrInsufficientFunds {
			return nil, ErrInsufficientFund
		}
		return nil, err
	}

	hold, _ := s.holdRepo.FindByReference(req.ReferenceID, req.ReferenceType)

	return &dto.HoldResponse{
		HoldID:        hold.ID,
		WalletID:      wallet.ID,
		Amount:        hold.Amount,
		ReferenceID:   hold.ReferenceID,
		ReferenceType: hold.ReferenceType,
		Status:        hold.Status,
	}, nil
}

func (s *WalletService) ReleaseHold(ctx context.Context, userID uint, req *dto.ReleaseHoldRequest) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	hold, err := s.holdRepo.FindByReference(req.ReferenceID, req.ReferenceType)
	if err != nil {
		return ErrHoldNotFound
	}

	if hold.WalletID != wallet.ID {
		return ErrInvalidUser
	}

	return s.holdRepo.AtomicReleaseHold(req.ReferenceID, req.ReferenceType)
}

func (s *WalletService) Debit(ctx context.Context, userID uint, req *dto.DebitRequest) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	if err := s.walletRepo.AtomicDebit(wallet.ID, req.Amount, req.ReferenceID, req.ReferenceType); err != nil {
		if err == repository.ErrInsufficientFunds {
			return ErrInsufficientFund
		}
		return err
	}

	if req.ReleaseHold {
		_ = s.holdRepo.AtomicReleaseHold(req.ReferenceID, req.ReferenceType)
	}

	return nil
}

func (s *WalletService) Credit(ctx context.Context, userID uint, req *dto.CreditRequest) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	return s.walletRepo.AtomicCredit(wallet.ID, req.Amount, req.ReferenceID, req.ReferenceType)
}

func (s *WalletService) Transfer(ctx context.Context, fromUserID, toUserID uint, amount float64, referenceID string) error {
	fromWallet, err := s.walletRepo.FindByUserID(fromUserID)
	if err != nil {
		return ErrWalletNotFound
	}

	toWallet, err := s.walletRepo.FindByUserID(toUserID)
	if err != nil {
		return ErrWalletNotFound
	}

	return s.walletRepo.AtomicTransfer(fromWallet.ID, toWallet.ID, amount, referenceID)
}

func (s *WalletService) TopUpStaff(ctx context.Context, mitraUserID, staffUserID uint, amount float64, referenceID string) error {
	mitraWallet, err := s.walletRepo.FindByUserID(mitraUserID)
	if err != nil {
		return ErrWalletNotFound
	}

	staffWallet, err := s.walletRepo.FindByUserID(staffUserID)
	if err != nil {
		return ErrWalletNotFound
	}

	if mitraWallet.Balance < amount {
		return ErrInsufficientFund
	}

	return s.walletRepo.AtomicTransfer(mitraWallet.ID, staffWallet.ID, amount, referenceID)
}

// TopUp credits the Mitra's own wallet
func (s *WalletService) TopUp(ctx context.Context, userID uint, amount float64) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}
	refID := "topup_" + uuid.New().String()
	return s.walletRepo.AtomicCredit(wallet.ID, amount, refID, "topup")
}

func (s *WalletService) GetEvents(ctx context.Context, userID uint, limit, offset int) ([]dto.TransactionResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	events, err := s.eventRepo.FindByWalletID(wallet.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TransactionResponse, len(events))
	for i, event := range events {
		responses[i] = dto.TransactionResponse{
			WalletID:      event.WalletID,
			EventType:     event.EventType,
			Amount:        event.Amount,
			BalanceBefore: event.BalanceBefore,
			BalanceAfter:  event.BalanceAfter,
			ReferenceID:   event.ReferenceID,
			CreatedAt:     event.CreatedAt.Unix(),
		}
	}

	return responses, nil
}

func (s *WalletService) Reconcile(ctx context.Context, userID uint) (*dto.ReconciliationResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrWalletNotFound
	}

	calculatedBalance, err := s.eventRepo.GetBalanceByEvents(wallet.ID)
	if err != nil {
		return nil, err
	}

	drift := wallet.Balance - calculatedBalance

	return &dto.ReconciliationResponse{
		WalletID:           wallet.ID,
		CurrentBalance:     wallet.Balance,
		CalculatedBalance:  calculatedBalance,
		Drift:              drift,
		IsBalanced:         drift == 0,
		ReconciledAt:       time.Now().Unix(),
	}, nil
}

func (s *WalletService) ValidateDailyLimit(ctx context.Context, userID uint, amount float64, limitType string, limitAmount float64) error {
	today := time.Now().Format("2006-01-02")

	key := fmt.Sprintf("daily_limit:%d:%s:%s", userID, limitType, today)
	existingAmount, err := s.redis.Get(ctx, key).Float64()
	if err != nil && err != redis.Nil {
		return err
	}

	if existingAmount+amount > limitAmount {
		return ErrDailyLimitExceeded
	}

	s.redis.IncrByFloat(ctx, key, amount)
	s.redis.Expire(ctx, key, 24*time.Hour)

	return nil
}

func (s *WalletService) CreateCommission(ctx context.Context, userID uint, amount float64, transactionID string, commissionType string, level int) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	referenceID := fmt.Sprintf("commission_%s_%d", transactionID, level)

	commission := &models.Commission{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Type:          commissionType,
		Level:         level,
	}

	if err := s.walletRepo.GetDB().Create(commission).Error; err != nil {
		return err
	}

	return s.walletRepo.AtomicCredit(wallet.ID, amount, referenceID, "commission")
}

func (s *WalletService) DebitForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	referenceID := fmt.Sprintf("debit_tx_%s", transactionID)
	referenceType := "transaction"

	if err := s.walletRepo.AtomicDebit(wallet.ID, amount, referenceID, referenceType); err != nil {
		return err
	}

	holdRef := fmt.Sprintf("hold_tx_%s", transactionID)
	hold, err := s.holdRepo.FindByReference(holdRef, referenceType)
	if err == nil && hold != nil {
		_ = s.holdRepo.AtomicReleaseHold(holdRef, referenceType)
	}

	return nil
}

func (s *WalletService) ReleaseHoldForTransaction(ctx context.Context, userID uint, transactionID string) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	referenceID := fmt.Sprintf("hold_tx_%s", transactionID)
	referenceType := "transaction"

	hold, err := s.holdRepo.FindByReference(referenceID, referenceType)
	if err != nil {
		return ErrHoldNotFound
	}

	if hold.WalletID != wallet.ID {
		return ErrInvalidUser
	}

	return s.holdRepo.AtomicReleaseHold(referenceID, referenceType)
}

func (s *WalletService) PlaceHoldForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return ErrWalletNotFound
	}

	referenceID := fmt.Sprintf("hold_tx_%s", transactionID)
	referenceType := "transaction"

	expiresAt := time.Now().Add(15 * time.Minute)

	return s.holdRepo.AtomicPlaceHold(wallet.ID, amount, referenceID, referenceType, &expiresAt)
}