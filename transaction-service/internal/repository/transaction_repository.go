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

func (r *TransactionRepository) GetProductByCode(ctx context.Context, code string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("code = ?", code).First(&product).Error
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

// Report aggregations

type SalesTrendResult struct {
	Date  string  `json:"date"`
	Sales float64 `json:"sales"`
	Count int     `json:"count"`
}

func (r *TransactionRepository) GetSalesTrend(startDate, endDate string) ([]SalesTrendResult, error) {
	var results []SalesTrendResult
	err := r.db.Model(&models.Transaction{}).
		Select("DATE(created_at) as date, SUM(amount) as sales, COUNT(*) as count").
		Where("status = ?", "success").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

type StaffPerformanceResult struct {
	StaffID          uint    `json:"staff_id"`
	StaffName        string  `json:"staff_name"`
	TransactionCount int     `json:"transaction_count"`
	TotalSales       float64 `json:"total_sales"`
	TotalCommission  float64 `json:"total_commission"`
	SuccessRate      float64 `json:"success_rate"`
}

func (r *TransactionRepository) GetStaffPerformance(startDate, endDate string) ([]StaffPerformanceResult, error) {
	var results []StaffPerformanceResult
	// This query joins with users to get staff name, and aggregates commission and sales
	err := r.db.Model(&models.Commission{}).
		Select("commissions.user_id as staff_id, users.full_name as staff_name, COUNT(*) as transaction_count, SUM(commissions.amount) as total_commission, SUM(case when commissions.status = 'paid' then 1 else 0 end) as paid_count").
		Joins("LEFT JOIN users ON commissions.user_id = users.id").
		Where("commissions.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("commissions.user_id, users.full_name").
		Order("total_commission DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	// Compute success rate from transactions separately? For now return.
	return results, err
}

type KPIResult struct {
	TotalSales     float64 `json:"total_sales"`
	PlatformProfit float64 `json:"platform_profit"`
	SuccessCount   int     `json:"success_count"`
	TotalCount     int     `json:"total_count"`
	StaffCount     int64   `json:"staff_count"`
}

func (r *TransactionRepository) GetKPIs(startDate, endDate string, userID uint) (*KPIResult, error) {
	var result KPIResult
	base := r.db.Model(&models.Transaction{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)
	if userID > 0 {
		base = base.Where("user_id = ?", userID)
	}

	// Total sales from successful transactions
	err := base.Select("COALESCE(SUM(amount),0) as total_sales").
		Where("status = ?", "success").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// Platform profit (sum of margin)
	err = base.Select("COALESCE(SUM(margin),0) as platform_profit").
		Where("status = ?", "success").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// Success count
	var successCount int64
	err = base.Where("status = ?", "success").Count(&successCount).Error
	if err != nil {
		return nil, err
	}
	result.SuccessCount = int(successCount)

	// Total count (all statuses)
	var total int64
	err = base.Count(&total).Error
	if err != nil {
		return nil, err
	}
	result.TotalCount = int(total)

	// Staff count (distinct user_ids)
	var staffCount int64
	err = base.Distinct("user_id").Count(&staffCount).Error
	if err != nil {
		return nil, err
	}
	result.StaffCount = staffCount

	return &result, nil
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