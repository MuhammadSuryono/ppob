package dto

import "time"

type ProductResponse struct {
	ID          uint      `json:"id"`
	ProductCode string    `json:"product_code"`
	ProductName string    `json:"product_name"`
	CategoryID  uint      `json:"category_id"`
	Provider    string    `json:"provider"`
	Price       float64   `json:"price"`
	PriceAPI    float64   `json:"price_api"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	ProductCode string  `json:"product_code" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	Provider    string  `json:"provider"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}

type UpdateProductRequest struct {
	ProductName string  `json:"product_name"`
	CategoryID  uint    `json:"category_id"`
	Provider    string  `json:"provider"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
}

type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64            `json:"total"`
	Limit    int              `json:"limit"`
	Offset   int              `json:"offset"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SyncResponse struct {
	Message   string `json:"message"`
	SyncedAt  string `json:"synced_at"`
	Count     int    `json:"count"`
}

type SyncStatusResponse struct {
	Prepaid  SyncTime  `json:"prepaid"`
	Postpaid SyncTime  `json:"postpaid"`
}

type SyncTime struct {
	LastSync int64 `json:"last_sync"`
}

type PriceValidationRequest struct {
	ProductCode  string             `json:"product_code" binding:"required"`
	SellingPrice float64            `json:"selling_price" binding:"required,gt=0"`
}

type PriceValidationResponse struct {
	Valid           bool    `json:"valid"`
	ProductCode     string  `json:"product_code"`
	SellingPrice    float64 `json:"selling_price"`
	PlatformPrice   float64 `json:"platform_price"`
	Margin          float64 `json:"margin"`
	MarginPercent   float64 `json:"margin_percent"`
	ValidationError string  `json:"validation_error,omitempty"`
}