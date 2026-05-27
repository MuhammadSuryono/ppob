package dto

type DigiflazzTransactionRequest struct {
	ProductCode    string  `json:"product_code" binding:"required"`
	CustomerNumber string  `json:"customer_number" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	RefID          string  `json:"ref_id" binding:"required"`
}

type DigiflazzTransactionResponse struct {
	Success      bool    `json:"success"`
	RefID        string  `json:"ref_id"`
	TrxID        string  `json:"trx_id"`
	Message      string  `json:"message"`
	Price        float64 `json:"price"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status"`
	ScCode       string  `json:"sc_code"`
	ScMessage    string  `json:"sc_message"`
}

type DigiflazzCallbackRequest struct {
	RefID     string  `json:"ref_id"`
	TrxID     string  `json:"trx_id"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
	Price     float64 `json:"price"`
	ScCode    string  `json:"sc_code"`
	ScMessage string  `json:"sc_message"`
}

type UpdateStatusRequest struct {
	Status         string `json:"status"`
	ProviderRef    string `json:"provider_ref"`
	ProviderStatus string `json:"provider_status"`
	Message        string `json:"message"`
}

type ProviderResponse struct {
	Provider   string `json:"provider"`
	Name       string `json:"name"`
	APIURL     string `json:"api_url"`
	IsActive   bool   `json:"is_active"`
	RateLimit  int    `json:"rate_limit"`
	TimeoutSec int    `json:"timeout_sec"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}