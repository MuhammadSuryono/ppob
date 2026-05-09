package dto

type BalanceResponse struct {
	WalletID          uint    `json:"wallet_id"`
	UserID            uint    `json:"user_id"`
	Balance           float64 `json:"balance"`
	HoldAmount        float64 `json:"hold_amount"`
	Available         float64 `json:"available"`
	CalculatedBalance float64 `json:"calculated_balance"`
	Currency          string  `json:"currency"`
}

type HoldRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID   string  `json:"reference_id" binding:"required"`
	ReferenceType string  `json:"reference_type" binding:"required"`
	ExpiresAt     int64   `json:"expires_at"`
}

type HoldResponse struct {
	HoldID        uint    `json:"hold_id"`
	WalletID      uint    `json:"wallet_id"`
	Amount        float64 `json:"amount"`
	ReferenceID   string  `json:"reference_id"`
	ReferenceType string  `json:"reference_type"`
	Status        string  `json:"status"`
}

type ReleaseHoldRequest struct {
	ReferenceID   string `json:"reference_id" binding:"required"`
	ReferenceType string `json:"reference_type" binding:"required"`
}

type DebitRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID   string  `json:"reference_id" binding:"required"`
	ReferenceType string  `json:"reference_type" binding:"required"`
	ReleaseHold   bool    `json:"release_hold"`
}

type CreditRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID   string  `json:"reference_id" binding:"required"`
	ReferenceType string  `json:"reference_type" binding:"required"`
}

type TransactionResponse struct {
	WalletID      uint    `json:"wallet_id"`
	EventType     string  `json:"event_type"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	ReferenceID   string  `json:"reference_id"`
	CreatedAt     int64   `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}