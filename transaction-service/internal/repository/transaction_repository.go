package repository

import (
	"context"

	"github.com/yontech/ppob/transaction-service/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) FindByTransactionID(transactionID string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Where("transaction_id = ?", transactionID).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *TransactionRepository) List(userID uint, status string, startDate, endDate string, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error

	return transactions, total, err
}

func (r *TransactionRepository) FindByProviderRef(providerRef string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Where("provider_ref = ?", providerRef).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) FindByCustomerNumber(customerNumber string, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("customer_number = ?", customerNumber).Order("created_at DESC").Limit(limit).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *TransactionRepository) GetProductByID(ctx context.Context, productID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *TransactionRepository) GetMitraProductPrice(ctx context.Context, mitraID uint, productCode string) (float64, error) {
	return 0, nil // Not implemented
}

func (r *TransactionRepository) GetActiveWallet(ctx context.Context, userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ? AND status = ?", userID, "active").First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *TransactionRepository) GetDailyLimit(ctx context.Context, userID uint, date string) (*models.DailyLimit, error) {
	var limit models.DailyLimit
	err := r.db.Where("user_id = ? AND date = ?", userID, date).First(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *TransactionRepository) GetIdempotencyKey(ctx context.Context, key string) (string, error) {
	return "", nil // Not implemented
}

type MarginRepository struct {
	db *gorm.DB
}

func NewMarginRepository(db *gorm.DB) *MarginRepository {
	return &MarginRepository{db: db}
}

func (r *MarginRepository) GetMarginByProductCode(productCode string) (float64, error) {
	var margin struct {
		Margin float64
	}
	err := r.db.Table("margin_settings").Where("product_code = ? AND is_active = ?", productCode, true).
		Select("margin").First(&margin).Error
	return margin.Margin, err
}

func (r *MarginRepository) GetMarkupByProductCode(productCode string) (float64, error) {
	var margin struct {
		Markup float64
	}
	err := r.db.Table("margin_settings").Where("product_code = ? AND is_active = ?", productCode, true).
		Select("markup").First(&margin).Error
	return margin.Markup, err
}

type DailyLimitRepository struct {
	db *gorm.DB
}

func NewDailyLimitRepository(db *gorm.DB) *DailyLimitRepository {
	return &DailyLimitRepository{db: db}
}

func (r *DailyLimitRepository) GetOrCreate(userID uint, date string) (*models.DailyLimit, error) {
	var limit models.DailyLimit
	err := r.db.Where("user_id = ? AND date = ?", userID, date).First(&limit).Error
	if err != nil {
		limit = models.DailyLimit{UserID: userID, Date: date, Count: 0, TotalAmount: 0}
		err = r.db.Create(&limit).Error
	}
	return &limit, err
}

func (r *DailyLimitRepository) Increment(userID uint, date string, amount float64) error {
	return r.db.Model(&models.DailyLimit{}).
		Where("user_id = ? AND date = ?", userID, date).
		Updates(map[string]interface{}{"count": gorm.Expr("count + 1"), "total_amount": gorm.Expr("total_amount + ?", amount)}).Error
}