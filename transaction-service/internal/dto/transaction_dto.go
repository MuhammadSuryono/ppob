package dto

import "time"

type CreateTransactionRequest struct {
	ProductCode     string  `json:"product_code" binding:"required"`
	CustomerNumber  string  `json:"customer_number" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	AuthorizeID     string  `json:"authorize_id" binding:"required"`
}

type TransactionResponse struct {
	ID               uint       `json:"id"`
	TransactionID    string     `json:"transaction_id"`
	UserID           uint       `json:"user_id"`
	ProductCode      string     `json:"product_code"`
	CustomerNumber   string     `json:"customer_number"`
	Amount           float64    `json:"amount"`
	Price            float64    `json:"price"`
	Margin           float64    `json:"margin"`
	Status           string     `json:"status"`
	ProviderRef      string     `json:"provider_ref"`
	Message          string     `json:"message"`
	CompletedAt      *time.Time `json:"completed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type UpdateStatusRequest struct {
	Status          string `json:"status" binding:"required"`
	ProviderRef     string `json:"provider_ref"`
	ProviderStatus  string `json:"provider_status"`
	Message         string `json:"message"`
}

type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                  `json:"total"`
	Limit        int                    `json:"limit"`
	Offset       int                    `json:"offset"`
}

type TransactionHistoryResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	NextCursor  string                 `json:"next_cursor"`
	Total       int64                  `json:"total"`
	HasMore     bool                   `json:"has_more"`
}

type InitiateTransactionRequest struct {
	ProductCode     string  `json:"product_code" binding:"required"`
	ProductID       uint    `json:"product_id"`
	CustomerNumber  string  `json:"customer_number" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	SellingPrice    float64 `json:"selling_price"`
	IdempotencyKey  string  `json:"idempotency_key"`
	AuthorizeID     string  `json:"authorize_id" binding:"required"`
}

type InitiateTransactionResponse struct {
	TransactionID   string               `json:"transaction_id"`
	Status          string               `json:"status"`
	Amount          float64              `json:"amount"`
	Price           float64              `json:"price"`
	CustomerNumber  string               `json:"customer_number"`
	ProductCode     string               `json:"product_code"`
	CreatedAt       int64                `json:"created_at"`
	Message         string               `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type FinancialIntegrityCheckResponse struct {
	TransactionID string `json:"transaction_id"`
	BalanceCheck  bool   `json:"balance_check"`
	MarginCheck   bool   `json:"margin_check"`
	IsHealthy     bool   `json:"is_healthy"`
	CheckedAt     int64  `json:"checked_at"`
}

type MarginCalculationRequest struct {
	ProductCode    string  `json:"product_code"`
	Amount         float64 `json:"amount"`
	SellingPrice   float64 `json:"selling_price"`
	StaffID        uint    `json:"staff_id"`
	StaffUserID    uint    `json:"staff_user_id"`
	ProductCost    float64 `json:"product_cost"`
}

type MarginCalculationResponse struct {
	SellingPrice      float64 `json:"selling_price"`
	ProductCost       float64 `json:"product_cost"`
	Margin            float64 `json:"margin"`
	MarginPercentage  float64 `json:"margin_percentage"`
	Commission        float64 `json:"commission"`
	ProductCode       string  `json:"product_code"`
	PlatformPrice     float64 `json:"platform_price"`
	StaffMargin       float64 `json:"staff_margin"`
	CommissionType    string  `json:"commission_type"`
}

type CancelTransactionRequest struct {
	Reason string `json:"reason"`
}

type DigiflazzWebhookRequest struct {
	RefID      string `json:"ref_id"`
	TrxID      string `json:"trx_id"`
	Status     string `json:"status"`
	SCCode     string `json:"sc_code"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

// Report DTOs

type ReportKPIResponse struct {
	TotalSales          float64 `json:"total_sales"`
	PlatformProfit      float64 `json:"platform_profit"`
	StaffCount          int     `json:"staff_count"`
	SuccessRate         float64 `json:"success_rate"`
	TransactionCount    int     `json:"transaction_count"`
	PeriodStart         string  `json:"period_start"`
	PeriodEnd           string  `json:"period_end"`
}

type ReportSalesTrendItem struct {
	Date  string  `json:"date"`
	Sales float64 `json:"sales"`
	Count int     `json:"count"`
}

type ReportStaffPerformanceItem struct {
	StaffID            uint    `json:"staff_id"`
	StaffName          string  `json:"staff_name"`
	TransactionCount   int     `json:"transaction_count"`
	TotalSales         float64 `json:"total_sales"`
	TotalCommission    float64 `json:"total_commission"`
	SuccessRate        float64 `json:"success_rate"`
}

type ReportsResponse struct {
	KPIs               []ReportKPIResponse               `json:"kpis"`
	SalesTrend        []ReportSalesTrendItem            `json:"sales_trend"`
	StaffPerformance  []ReportStaffPerformanceItem      `json:"staff_performance"`
}