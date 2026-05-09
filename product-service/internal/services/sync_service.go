package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/product-service/config"
	"github.com/yontech/ppob/product-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrSyncFailed      = errors.New("product sync failed")
	ErrLockAcquisition = errors.New("failed to acquire sync lock")
)

type ProductSyncService struct {
	db       *gorm.DB
	redis    *redis.Client
	cfg      *config.Config
}

func NewProductSyncService(db *gorm.DB, redis *redis.Client, cfg *config.Config) *ProductSyncService {
	return &ProductSyncService{db: db, redis: redis, cfg: cfg}
}

type DigiflazzProduct struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Provider     string  `json:"provider"`
	Status       string  `json:"status"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
}

type DigiflazzPriceListResponse struct {
	Success bool               `json:"success"`
	Data    []DigiflazzProduct `json:"data"`
	Message string             `json:"message"`
}

func (s *ProductSyncService) SyncPrepaidProducts(ctx context.Context) error {
	lockKey := "product_sync:prepaid:lock"
	lockAcquired, err := s.acquireLock(ctx, lockKey, 5*time.Minute)
	if err != nil {
		return err
	}
	if !lockAcquired {
		return ErrLockAcquisition
	}
	defer s.releaseLock(ctx, lockKey)

	products, err := s.fetchDigiflazzProducts(ctx, "prepaid")
	if err != nil {
		return err
	}

	now := time.Now()
	for _, p := range products {
		categoryID, _ := s.getOrCreateCategory(ctx, p.Category)

		var existing models.Product
		err := s.db.Where("product_code = ?", p.Code).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			product := models.Product{
				ProductCode: p.Code,
				ProductName: p.Name,
				CategoryID:  categoryID,
				Provider:    p.Provider,
				Price:       p.Price,
				PriceAPI:    p.Price,
				Stock:       -1,
				Status:      s.mapStatus(p.Status),
				Description: p.Description,
				LastSyncAt:  &now,
			}
			s.db.Create(&product)
		} else if err == nil {
			existing.ProductName = p.Name
			existing.CategoryID = categoryID
			existing.Provider = p.Provider
			existing.PriceAPI = p.Price
			existing.LastSyncAt = &now
			if existing.Price == existing.PriceAPI {
				existing.Price = p.Price
			}
			s.db.Save(&existing)
		}
	}

	s.setLastSyncTime(ctx, "prepaid")
	return nil
}

func (s *ProductSyncService) SyncPostpaidProducts(ctx context.Context) error {
	lockKey := "product_sync:postpaid:lock"
	lockAcquired, err := s.acquireLock(ctx, lockKey, 2*time.Minute)
	if err != nil {
		return err
	}
	if !lockAcquired {
		return ErrLockAcquisition
	}
	defer s.releaseLock(ctx, lockKey)

	products, err := s.fetchDigiflazzProducts(ctx, "postpaid")
	if err != nil {
		return err
	}

	now := time.Now()
	for _, p := range products {
		categoryID, _ := s.getOrCreateCategory(ctx, p.Category)

		var existing models.Product
		err := s.db.Where("product_code = ?", p.Code).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			product := models.Product{
				ProductCode: p.Code,
				ProductName: p.Name,
				CategoryID:  categoryID,
				Provider:    p.Provider,
				Price:       p.Price,
				PriceAPI:    p.Price,
				Stock:       -1,
				Status:      s.mapStatus(p.Status),
				Description: p.Description,
				LastSyncAt:  &now,
			}
			s.db.Create(&product)
		} else if err == nil {
			existing.ProductName = p.Name
			existing.CategoryID = categoryID
			existing.Provider = p.Provider
			existing.PriceAPI = p.Price
			existing.LastSyncAt = &now
			if existing.Price == existing.PriceAPI {
				existing.Price = p.Price
			}
			s.db.Save(&existing)
		}
	}

	s.setLastSyncTime(ctx, "postpaid")
	return nil
}

func (s *ProductSyncService) fetchDigiflazzProducts(ctx context.Context, productType string) ([]DigiflazzProduct, error) {
	if s.cfg.DigiflazzKey == "" {
		return s.getMockProducts(productType), nil
	}

	url := fmt.Sprintf("%s/pricelist", s.cfg.DigiflazzURL)
	payload := map[string]interface{}{
		"username": s.cfg.DigiflazzKey,
		"sign":     fmt.Sprintf("%s%s", s.cfg.DigiflazzKey, s.cfg.DigiflazzSecret),
		"cmd":      productType,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DigiflazzPriceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("digiflazz error: %s", result.Message)
	}

	return result.Data, nil
}

func (s *ProductSyncService) getMockProducts(productType string) []DigiflazzProduct {
	if productType == "prepaid" {
		return []DigiflazzProduct{
			{Code: "PREPAID_SIMPATIS_5K", Name: "Simpati 5.000", Price: 5500, Provider: "Telkomsel", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_SIMPATIS_10K", Name: "Simpati 10.000", Price: 10500, Provider: "Telkomsel", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_SIMPATIS_20K", Name: "Simpati 20.000", Price: 20500, Provider: "Telkomsel", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_AS_5K", Name: "AS 5.000", Price: 5500, Provider: "Telkomsel", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_XL_10K", Name: "XL 10.000", Price: 11000, Provider: "XL", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_AXIS_10K", Name: "Axis 10.000", Price: 11000, Provider: "Axis", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_THREE_5K", Name: "Three 5.000", Price: 5500, Provider: "Three", Status: "active", Category: "Pulsa"},
			{Code: "PREPAID_TOKEN_LISTRIK_20K", Name: "Token Listrik 20.000", Price: 21000, Provider: "PLN", Status: "active", Category: "Token Listrik"},
			{Code: "PREPAID_TOKEN_LISTRIK_50K", Name: "Token Listrik 50.000", Price: 51000, Provider: "PLN", Status: "active", Category: "Token Listrik"},
			{Code: "PREPAID_TOKEN_LISTRIK_100K", Name: "Token Listrik 100.000", Price: 101000, Provider: "PLN", Status: "active", Category: "Token Listrik"},
		}
	}
	return []DigiflazzProduct{
		{Code: "POSTPAID_TAGIHAN_TELKOM", Name: "Tagihan Telkomsel", Price: 0, Provider: "Telkomsel", Status: "active", Category: "Tagihan"},
		{Code: "POSTPAID_TAGIHAN_XL", Name: "Tagihan XL", Price: 0, Provider: "XL", Status: "active", Category: "Tagihan"},
		{Code: "POSTPAID_TAGIHAN_MOTOR", Name: "BPJS Kendaraan", Price: 0, Provider: "Samsung", Status: "active", Category: "Finance"},
	}
}

func (s *ProductSyncService) getOrCreateCategory(ctx context.Context, name string) (uint, error) {
	if name == "" {
		name = "Uncategorized"
	}

	var category models.Category
	err := s.db.Where("name = ?", name).First(&category).Error
	if err == gorm.ErrRecordNotFound {
		category = models.Category{
			Name:       name,
			Code:       fmt.Sprintf("CAT_%d", time.Now().Unix()),
			Status:     "active",
			SortOrder:  0,
		}
		s.db.Create(&category)
	}
	return category.ID, err
}

func (s *ProductSyncService) mapStatus(status string) string {
	switch status {
	case "active", "Active", "ON":
		return "active"
	case "maintenance", "Maintenance", "OFF":
		return "maintenance"
	default:
		return "inactive"
	}
}

func (s *ProductSyncService) acquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	result, err := s.redis.SetNX(ctx, key, "locked", ttl).Result()
	return result, err
}

func (s *ProductSyncService) releaseLock(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key).Err()
}

func (s *ProductSyncService) setLastSyncTime(ctx context.Context, productType string) error {
	key := fmt.Sprintf("product_sync:last:%s", productType)
	return s.redis.Set(ctx, key, time.Now().Unix(), 24*time.Hour).Err()
}

func (s *ProductSyncService) GetLastSyncTime(ctx context.Context, productType string) (time.Time, error) {
	key := fmt.Sprintf("product_sync:last:%s", productType)
	ts, err := s.redis.Get(ctx, key).Int64()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ts, 0), nil
}

type PriceValidationService struct {
	db *gorm.DB
}

func NewPriceValidationService(db *gorm.DB) *PriceValidationService {
	return &PriceValidationService{db: db}
}

type PriceValidationResult struct {
	Valid           bool    `json:"valid"`
	ProductCode     string  `json:"product_code"`
	SellingPrice    float64 `json:"selling_price"`
	PlatformPrice   float64 `json:"platform_price"`
	Margin          float64 `json:"margin"`
	MarginPercent   float64 `json:"margin_percent"`
	ValidationError string  `json:"validation_error,omitempty"`
}

func (s *PriceValidationService) ValidatePrice(productCode string, sellingPrice float64) *PriceValidationResult {
	var product models.Product
	err := s.db.Where("product_code = ?", productCode).First(&product).Error
	if err != nil {
		return &PriceValidationResult{
			Valid:           false,
			ProductCode:     productCode,
			SellingPrice:    sellingPrice,
			ValidationError: "Product not found",
		}
	}

	platformPrice := product.PriceAPI
	if platformPrice <= 0 {
		platformPrice = product.Price
	}

	margin := sellingPrice - platformPrice

	result := &PriceValidationResult{
		Valid:           true,
		ProductCode:     productCode,
		SellingPrice:    sellingPrice,
		PlatformPrice:   platformPrice,
		Margin:          margin,
	}

	if margin < 0 {
		result.Valid = false
		result.ValidationError = fmt.Sprintf("Selling price (%.0f) cannot be below platform price (%.0f)", sellingPrice, platformPrice)
	}

	if platformPrice > 0 {
		result.MarginPercent = (margin / platformPrice) * 100
	}

	return result
}

func (s *PriceValidationService) BatchValidate(pricing map[string]float64) map[string]*PriceValidationResult {
	results := make(map[string]*PriceValidationResult)
	for productCode, price := range pricing {
		results[productCode] = s.ValidatePrice(productCode, price)
	}
	return results
}