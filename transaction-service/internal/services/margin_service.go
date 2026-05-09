package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrStaffNotFound       = errors.New("staff not found")
	ErrInvalidPrice        = errors.New("invalid price calculation")
)

type MarginService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewMarginService(db *gorm.DB, cfg *config.Config) *MarginService {
	return &MarginService{db: db, cfg: cfg}
}

type MarginSettings struct {
	ID                  uint    `gorm:"primaryKey"`
	UserID              uint    `gorm:"uniqueIndex"`
	CommissionType      string  `gorm:"size:20;default:MarginShare"`
	GlobalMarginPercent float64 `gorm:"default:0"`
	FixedAllowance     float64 `gorm:"default:0"`
	IsActive            bool    `gorm:"default:true"`
}

type ProductMarginOverride struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"index"`
	ProductCode  string  `gorm:"size:100"`
	MarginPercent float64 `gorm:"default:0"`
	FixedMargin  float64 `gorm:"default:0"`
	IsActive     bool    `gorm:"default:true"`
}

func (s *MarginService) CalculateMargin(ctx context.Context, req *dto.MarginCalculationRequest) (*dto.MarginCalculationResponse, error) {
	product, err := s.getProductByCode(req.ProductCode)
	if err != nil {
		return nil, ErrProductNotFound
	}

	platformPrice := product.PriceAPI
	if platformPrice <= 0 {
		platformPrice = product.Price
	}

	margin := req.SellingPrice - platformPrice
	if margin < 0 {
		margin = 0
	}

	staffSettings, err := s.getStaffMarginSettings(req.StaffUserID)
	if err != nil {
		return nil, err
	}

	override, _ := s.getProductOverride(req.StaffUserID, req.ProductCode)

	var staffMargin, commission float64
	var commissionType string

	if staffSettings.CommissionType == "FixedAllowance" {
		commissionType = "FixedAllowance"
		staffMargin = staffSettings.FixedAllowance
		commission = staffSettings.FixedAllowance
	} else {
		commissionType = "MarginShare"
		marginPercent := staffSettings.GlobalMarginPercent
		if override != nil && override.IsActive {
			marginPercent = override.MarginPercent
		}

		staffMargin = margin * (marginPercent / 100)
		commission = staffMargin
	}

	return &dto.MarginCalculationResponse{
		ProductCode:     req.ProductCode,
		SellingPrice:    req.SellingPrice,
		PlatformPrice:   platformPrice,
		Margin:          margin,
		StaffMargin:     staffMargin,
		Commission:      commission,
		CommissionType:  commissionType,
	}, nil
}

func (s *MarginService) getProductByCode(productCode string) (*models.Product, error) {
	var product models.Product
	err := s.db.Where("product_code = ?", productCode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *MarginService) getStaffMarginSettings(userID uint) (*MarginSettings, error) {
	var settings MarginSettings
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &MarginSettings{
				UserID:              userID,
				CommissionType:      "MarginShare",
				GlobalMarginPercent: 10.0,
				FixedAllowance:      0,
				IsActive:            true,
			}, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (s *MarginService) getProductOverride(userID uint, productCode string) (*ProductMarginOverride, error) {
	var override ProductMarginOverride
	err := s.db.Where("user_id = ? AND product_code = ? AND is_active = ?", userID, productCode, true).First(&override).Error
	if err != nil {
		return nil, err
	}
	return &override, nil
}

func (s *MarginService) CalculateTransactionMargin(userID uint, productCode string, sellingPrice float64) (margin, staffCommission float64, err error) {
	marginResp, err := s.CalculateMargin(context.Background(), &dto.MarginCalculationRequest{
		ProductCode:  productCode,
		SellingPrice: sellingPrice,
		StaffUserID:  userID,
	})
	if err != nil {
		return 0, 0, err
	}

	return marginResp.Margin, marginResp.Commission, nil
}

type WalletClient interface {
	Credit(amount float64, referenceID, referenceType string) error
}

type CommissionService struct {
	db          *gorm.DB
	marginSvc   *MarginService
	walletSvc   WalletClient
}

func NewCommissionService(db *gorm.DB, marginSvc *MarginService, walletSvc WalletClient) *CommissionService {
	return &CommissionService{
		db:        db,
		marginSvc: marginSvc,
		walletSvc: walletSvc,
	}
}

func (s *CommissionService) CreateAndCreditCommission(ctx context.Context, userID uint, transactionID string, amount float64, commissionType string, level int) (*models.Commission, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	txID, _ := strconv.ParseUint(transactionID, 10, 32)
	commission := &models.Commission{
		UserID:        userID,
		TransactionID: uint(txID),
		Amount:        amount,
		Type:          commissionType,
		Level:         level,
		Status:        "completed",
		EarnedAt:      &now,
	}

	if err := tx.Create(commission).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if s.walletSvc != nil {
		referenceID := "commission_" + transactionID
		if err := s.walletSvc.Credit(amount, referenceID, "commission"); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return commission, nil
}

func (s *CommissionService) CreateCommission(ctx context.Context, userID uint, transactionID string, amount float64, commissionType string, level int) (*models.Commission, error) {
	txID, _ := strconv.ParseUint(transactionID, 10, 32)
	now := time.Now()
	commission := &models.Commission{
		UserID:        userID,
		TransactionID: uint(txID),
		Amount:        amount,
		Type:          commissionType,
		Level:         level,
		Status:        "completed",
		EarnedAt:      &now,
	}

	if err := s.db.Create(commission).Error; err != nil {
		return nil, err
	}

	return commission, nil
}

func (s *CommissionService) GetCommissionByTransaction(transactionID string) ([]models.Commission, error) {
	var commissions []models.Commission
	err := s.db.Where("transaction_id = ?", transactionID).Find(&commissions).Error
	return commissions, err
}

func (s *CommissionService) GetTotalCommissionByUser(userID uint) (float64, error) {
	var result struct {
		Total float64
	}
	err := s.db.Model(&models.Commission{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&result).Error
	return result.Total, err
}