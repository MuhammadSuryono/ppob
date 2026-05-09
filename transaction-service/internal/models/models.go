package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"size:255;uniqueIndex" json:"email"`
	Phone     string         `gorm:"size:20;uniqueIndex" json:"phone"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Wallet struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"index;not null" json:"user_id"`
	RoleID             uint           `gorm:"index" json:"role_id"`
	Balance            float64        `gorm:"default:0" json:"balance"`
	AvailableBalance   float64        `gorm:"default:0" json:"available_balance"`
	BalanceAvailable  float64        `gorm:"default:0" json:"balance_available"`
	HoldAmount        float64        `gorm:"default:0" json:"hold_amount"`
	Status            string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type Transaction struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	TransactionID    string         `gorm:"uniqueIndex;size:100;not null" json:"transaction_id"`
	UserID            uint           `gorm:"index;not null" json:"user_id"`
	WalletID          *uint          `gorm:"index" json:"wallet_id"`
	ProductID         uint           `gorm:"index;not null" json:"product_id"`
	ProductCode       string         `gorm:"size:100;not null" json:"product_code"`
	CustomerNumber    string         `gorm:"size:50;not null" json:"customer_number"`
	Amount            float64        `gorm:"not null" json:"amount"`
	Price             float64        `gorm:"not null" json:"price"`
	SellingPrice      float64        `gorm:"default:0" json:"selling_price"`
	Margin            float64        `gorm:"default:0" json:"margin"`
	Status            string         `gorm:"size:20;default:pending" json:"status"`
	ProviderRef       string         `gorm:"size:100" json:"provider_ref"`
	ProviderStatus    string         `gorm:"size:20" json:"provider_status"`
	Message           string         `gorm:"type:text" json:"message"`
	CompletedAt       *time.Time     `json:"completed_at"`
	PreviousStatus    string         `gorm:"size:50" json:"previous_status"`
	StatusChangeReason string        `gorm:"size:255" json:"status_change_reason"`
	ReconciledAt      *time.Time     `json:"reconciled_at"`
	HoldReleasedAt    *time.Time     `json:"hold_released_at"`
	DigiflazzRC       string         `gorm:"size:10" json:"digiflazz_rc"`
	WebhookReceivedAt *time.Time     `json:"webhook_received_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type DailyLimit struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Date         string    `gorm:"size:10;not null" json:"date"`
	Count        int       `gorm:"default:0" json:"count"`
	MaxCount     int       `gorm:"default:100" json:"max_count"`
	TotalAmount  float64   `gorm:"default:0" json:"total_amount"`
	MaxAmount    float64   `gorm:"default:10000000" json:"max_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Commission struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TransactionID  uint           `gorm:"index;not null" json:"transaction_id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	WalletID       *uint          `gorm:"index" json:"wallet_id"`
	Amount         float64        `gorm:"not null" json:"amount"`
	Type           string         `gorm:"size:20;not null" json:"type"`
	Level          int            `gorm:"default:1" json:"level"`
	Status         string         `gorm:"size:20;default:pending" json:"status"`
	StaffID        uint           `gorm:"index" json:"staff_id"`
	SchemeUsed     string         `gorm:"size:50" json:"scheme_used"`
	MarginAmount   float64        `json:"margin_amount"`
	PaidAt         *time.Time     `json:"paid_at"`
	EarnedAt       *time.Time     `json:"earned_at"`
	CreatedAt      time.Time      `json:"created_at"`
}

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Category      string    `gorm:"size:50" json:"category"`
	Price         float64   `gorm:"not null" json:"price"`
	PriceAPI      float64   `gorm:"default:0" json:"price_api"`
	PlatformPrice float64  `gorm:"default:0" json:"platform_price"`
	Status        string    `gorm:"size:20;default:active" json:"status"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	Provider      string    `gorm:"size:50" json:"provider"`
	ProductType   string    `gorm:"size:20" json:"product_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PostpaidInquiry struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	InquiryID         string         `gorm:"uniqueIndex;size:100" json:"inquiry_id"`
	RefID             string         `gorm:"size:100" json:"ref_id"`
	ProductID         uint           `gorm:"index" json:"product_id"`
	CustomerNumber    string         `gorm:"size:50;not null" json:"customer_number"`
	CustomerNo        string         `gorm:"size:50" json:"customer_no"`
	CustomerName      string         `gorm:"size:255" json:"customer_name"`
	BillAmount        float64        `json:"bill_amount"`
	AdminFee          float64        `json:"admin_fee"`
	AdminAmount       float64        `json:"admin_amount"`
	TotalAmount       float64        `json:"total_amount"`
	BillDetails       map[string]interface{} `gorm:"type:jsonb" json:"bill_details"`
	SellingPrice      float64        `json:"selling_price"`
	Period            string         `gorm:"size:10" json:"period"`
	Status            string         `gorm:"size:20;default:active" json:"status"`
	ExpiresAt         time.Time      `json:"expires_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}