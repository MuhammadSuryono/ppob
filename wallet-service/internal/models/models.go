package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance     float64        `gorm:"default:0" json:"balance"`
	HoldAmount  float64        `gorm:"default:0" json:"hold_amount"`
	Currency    string         `gorm:"size:10;default:IDR" json:"currency"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WalletEvent struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	WalletID      uint           `gorm:"index;not null" json:"wallet_id"`
	EventType     string         `gorm:"size:50;not null" json:"event_type"`
	Amount        float64        `json:"amount"`
	BalanceBefore float64        `json:"balance_before"`
	BalanceAfter  float64        `json:"balance_after"`
	ReferenceID   string         `gorm:"size:100" json:"reference_id"`
	ReferenceType string         `gorm:"size:50" json:"reference_type"`
	Metadata      string         `gorm:"type:text" json:"metadata"`
	CreatedAt     time.Time      `json:"created_at"`
}

type Hold struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	WalletID      uint           `gorm:"index;not null" json:"wallet_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	ReferenceID   string         `gorm:"size:100" json:"reference_id"`
	ReferenceType string         `gorm:"size:50" json:"reference_type"`
	Status        string         `gorm:"size:20;default:active" json:"status"`
	ExpiresAt     *time.Time     `json:"expires_at"`
	CreatedAt     time.Time      `json:"created_at"`
	ReleasedAt    *time.Time     `json:"released_at"`
}

type Commission struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TransactionID string         `gorm:"index;size:100;not null" json:"transaction_id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	Type          string         `gorm:"size:20;not null" json:"type"`
	Level         int            `gorm:"default:1" json:"level"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	ReferenceID   string         `gorm:"size:100" json:"reference_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type DailyLimit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Date        string    `gorm:"size:10;not null" json:"date"`
	Count       int       `gorm:"default:0" json:"count"`
	TotalAmount float64   `gorm:"default:0" json:"total_amount"`
	MaxCount    int       `gorm:"default:0" json:"max_count"`
	MaxAmount   float64   `gorm:"default:0" json:"max_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}