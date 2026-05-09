package services

import (
	"context"
	"errors"
	"time"

	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/repository"
)

var (
	ErrValidationProductNotFound = errors.New("product not found or inactive")
	ErrPriceBelowCost            = errors.New("selling price below minimum allowed")
	ErrInsufficientBalance       = errors.New("insufficient wallet balance")
	ErrValidationDailyLimit      = errors.New("daily limit exceeded")
	ErrUserInactive              = errors.New("user account is inactive")
	ErrInvalidCustomerNo         = errors.New("invalid customer number")
	ErrIdempotencyConflict       = errors.New("duplicate request - idempotency key conflict")
)

type TransactionValidationService struct {
	transactionRepo *repository.TransactionRepository
	marginRepo      *repository.MarginRepository
	cfg             *config.Config
}

func NewTransactionValidationService(
	transactionRepo *repository.TransactionRepository,
	marginRepo *repository.MarginRepository,
	cfg *config.Config,
) *TransactionValidationService {
	return &TransactionValidationService{
		transactionRepo: transactionRepo,
		marginRepo:      marginRepo,
		cfg:             cfg,
	}
}

type ValidationResult struct {
	Valid               bool
	ActualSellingPrice  float64
	PlatformPrice       float64
	AvailableBalance    float64
	DailyCount          int64
	DailyAmount         float64
	MaxDailyCount       int64
	MaxDailyAmount      float64
	Error               error
}

func (s *TransactionValidationService) ValidateAndCreate(ctx context.Context, req dto.InitiateTransactionRequest, userID uint) (*ValidationResult, error) {
	result := &ValidationResult{Valid: false}

	user, err := s.transactionRepo.GetUserByID(ctx, userID)
	if err != nil || !user.IsActive {
		result.Error = ErrUserInactive
		return result, result.Error
	}

	product, err := s.transactionRepo.GetProductByID(ctx, req.ProductID)
	if err != nil || !product.IsActive {
		result.Error = ErrValidationProductNotFound
		return result, result.Error
	}

	mitraID := user.ID
	actualSellingPrice := product.PlatformPrice

	mitraPrice, err := s.transactionRepo.GetMitraProductPrice(ctx, mitraID, req.ProductCode)
	if err == nil && mitraPrice > 0 {
		actualSellingPrice = mitraPrice
	}

	if req.SellingPrice < actualSellingPrice {
		result.Error = ErrPriceBelowCost
		result.ActualSellingPrice = actualSellingPrice
		result.PlatformPrice = product.PlatformPrice
		return result, result.Error
	}

	wallet, err := s.transactionRepo.GetActiveWallet(ctx, userID)
	if err != nil {
		result.Error = err
		return result, result.Error
	}

	result.AvailableBalance = wallet.BalanceAvailable

	if wallet.BalanceAvailable < req.SellingPrice {
		result.Error = ErrInsufficientBalance
		return result, result.Error
	}

	dailyLimit, err := s.transactionRepo.GetDailyLimit(ctx, userID, time.Now().Format("2006-01-02"))
	if err == nil && dailyLimit != nil {
		result.DailyCount = int64(dailyLimit.Count)
		result.DailyAmount = dailyLimit.TotalAmount
		result.MaxDailyCount = int64(dailyLimit.MaxCount)
		result.MaxDailyAmount = dailyLimit.MaxAmount

		if dailyLimit.Count >= dailyLimit.MaxCount {
			result.Error = ErrValidationDailyLimit
			return result, result.Error
		}

		if dailyLimit.TotalAmount+req.SellingPrice > dailyLimit.MaxAmount {
			result.Error = ErrValidationDailyLimit
			return result, result.Error
		}
	}

	idempotencyKey := req.IdempotencyKey
	if idempotencyKey != "" {
		existing, err := s.transactionRepo.GetIdempotencyKey(ctx, idempotencyKey)
		if err == nil && existing != "" {
			result.Error = ErrIdempotencyConflict
			return result, result.Error
		}
	}

	result.Valid = true
	result.ActualSellingPrice = actualSellingPrice
	result.PlatformPrice = product.PlatformPrice

	return result, nil
}

func (s *TransactionValidationService) ValidateCustomerNo(productCode, customerNo string) error {
	minLen := 4
	maxLen := 25

	cleaned := customerNo
	if len(cleaned) < minLen || len(cleaned) > maxLen {
		return ErrInvalidCustomerNo
	}

	return nil
}