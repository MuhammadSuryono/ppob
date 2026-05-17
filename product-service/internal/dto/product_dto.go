package dto

import "time"

type ProductResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
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
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Brand       string  `json:"brand"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	Provider    string  `json:"provider"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
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
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	Icon            string `json:"icon"`
	SortOrder       int    `json:"sort_order"`
	Status          string `json:"status"`
	InputType       string `json:"input_type"`
	InputLabel      string `json:"input_label"`
	Placeholder     string `json:"placeholder"`
	ValidationRegex string `json:"validation_regex"`
}

type CreateCategoryRequest struct {
	Name            string `json:"name" binding:"required"`
	Code            string `json:"code" binding:"required"`
	Description     string `json:"description"`
	Icon            string `json:"icon"`
	SortOrder       int    `json:"sort_order"`
	InputType       string `json:"input_type"`
	InputLabel      string `json:"input_label"`
	Placeholder     string `json:"placeholder"`
	ValidationRegex string `json:"validation_regex"`
}

type UpdateCategoryRequest struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	Icon            string `json:"icon"`
	SortOrder       int    `json:"sort_order"`
	Status          string `json:"status"`
	InputType       string `json:"input_type"`
	InputLabel      string `json:"input_label"`
	Placeholder     string `json:"placeholder"`
	ValidationRegex string `json:"validation_regex"`
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
	Code         string             `json:"code" binding:"required"`
	SellingPrice float64            `json:"selling_price" binding:"required,gt=0"`
}

type PriceValidationResponse struct {
	Valid           bool    `json:"valid"`
	Code            string  `json:"code"`
	SellingPrice    float64 `json:"selling_price"`
	PlatformPrice   float64 `json:"platform_price"`
	Margin          float64 `json:"margin"`
	MarginPercent   float64 `json:"margin_percent"`
	ValidationError string  `json:"validation_error,omitempty"`
}