package dto

type ReconciliationResponse struct {
	WalletID          uint    `json:"wallet_id"`
	CurrentBalance    float64 `json:"current_balance"`
	CalculatedBalance float64 `json:"calculated_balance"`
	Drift             float64 `json:"drift"`
	IsBalanced        bool    `json:"is_balanced"`
	ReconciledAt      int64   `json:"reconciled_at"`
}

type TransferRequest struct {
	ToUserID   uint    `json:"to_user_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID string  `json:"reference_id" binding:"required"`
}

type TopUpStaffRequest struct {
	StaffUserID  uint    `json:"staff_user_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	ReferenceID  string  `json:"reference_id" binding:"required"`
}

type DailyLimitRequest struct {
	MaxCount  int     `json:"max_count"`
	MaxAmount float64 `json:"max_amount"`
}

type DailyLimitResponse struct {
	UserID     uint    `json:"user_id"`
	Date       string  `json:"date"`
	Count      int     `json:"count"`
	TotalAmount float64 `json:"total_amount"`
	MaxCount   int     `json:"max_count"`
	MaxAmount  float64 `json:"max_amount"`
	Remaining  int     `json:"remaining"`
}

type CommissionRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	TransactionID   string  `json:"transaction_id" binding:"required"`
	CommissionType  string  `json:"commission_type" binding:"required"`
	Level           int     `json:"level" binding:"required,min=1"`
}

type CommissionResponse struct {
	ID             uint    `json:"id"`
	TransactionID  string  `json:"transaction_id"`
	UserID         uint    `json:"user_id"`
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	Level          int     `json:"level"`
	Status         string  `json:"status"`
	CreatedAt      int64   `json:"created_at"`
}

type MarginCalculationRequest struct {
	ProductCode     string  `json:"product_code" binding:"required"`
	SellingPrice    float64 `json:"selling_price" binding:"required,gt=0"`
	StaffUserID     uint    `json:"staff_user_id" binding:"required"`
}

type MarginCalculationResponse struct {
	ProductCode       string  `json:"product_code"`
	SellingPrice      float64 `json:"selling_price"`
	PlatformPrice     float64 `json:"platform_price"`
	Margin            float64 `json:"margin"`
	StaffMargin       float64 `json:"staff_margin"`
	Commission        float64 `json:"commission"`
	CommissionType    string  `json:"commission_type"`
}

type FinancialIntegrityCheckResponse struct {
	WalletID           uint    `json:"wallet_id"`
	BalanceCheck       bool    `json:"balance_check"`
	HoldAmountCheck    bool    `json:"hold_amount_check"`
	EventCountCheck    bool    `json:"event_count_check"`
	DriftAmount        float64 `json:"drift_amount"`
	IsHealthy          bool    `json:"is_healthy"`
	CheckedAt          int64   `json:"checked_at"`
}