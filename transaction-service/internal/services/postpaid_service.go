package services

import (
	"errors"
	"time"

	"github.com/yontech/ppob/transaction-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrInquiryNotFound   = errors.New("inquiry not found")
	ErrInquiryExpired   = errors.New("inquiry has expired, please re-inquire")
	ErrInquiryNotMatching = errors.New("customer number does not match inquiry")
)

type PostpaidInquiryService struct {
	db *gorm.DB
}

func NewPostpaidInquiryService(db *gorm.DB) *PostpaidInquiryService {
	return &PostpaidInquiryService{db: db}
}

func (s *PostpaidInquiryService) CreateInquiry(req *PostpaidInquiryRequest) (*models.PostpaidInquiry, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	inquiry := &models.PostpaidInquiry{
		RefID:         req.RefID,
		ProductID:     req.ProductID,
		CustomerNo:    req.CustomerNo,
		CustomerName:  req.CustomerName,
		BillDetails:   req.BillDetails,
		AdminAmount:   req.AdminAmount,
		TotalAmount:   req.TotalAmount,
		SellingPrice:  req.SellingPrice,
		ExpiresAt:     expiresAt,
	}

	err := s.db.Create(inquiry).Error
	return inquiry, err
}

func (s *PostpaidInquiryService) GetInquiry(refID string) (*models.PostpaidInquiry, error) {
	var inquiry models.PostpaidInquiry
	err := s.db.Where("ref_id = ?", refID).First(&inquiry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInquiryNotFound
		}
		return nil, err
	}
	return &inquiry, nil
}

func (s *PostpaidInquiryService) ValidateForPayment(refID, customerNo string) (*models.PostpaidInquiry, error) {
	inquiry, err := s.GetInquiry(refID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(inquiry.ExpiresAt) {
		return nil, ErrInquiryExpired
	}

	if inquiry.CustomerNo != customerNo {
		return nil, ErrInquiryNotMatching
	}

	return inquiry, nil
}

func (s *PostpaidInquiryService) CleanupExpired() (int64, error) {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.PostpaidInquiry{})
	return result.RowsAffected, result.Error
}

type PostpaidInquiryRequest struct {
	RefID        string                 `json:"ref_id" binding:"required"`
	ProductID    uint                   `json:"product_id" binding:"required"`
	CustomerNo   string                 `json:"customer_no" binding:"required"`
	CustomerName string                 `json:"customer_name"`
	BillDetails  map[string]interface{} `json:"bill_details" binding:"required"`
	AdminAmount  float64                `json:"admin_amount" binding:"required"`
	TotalAmount  float64                `json:"total_amount" binding:"required"`
	SellingPrice float64                `json:"selling_price" binding:"required"`
}