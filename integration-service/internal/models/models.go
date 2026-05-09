package models

import (
	"time"

	"gorm.io/gorm"
)

type IntegrationLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Provider      string         `gorm:"size:50;not null" json:"provider"`
	Action        string         `gorm:"size:50;not null" json:"action"`
	RequestID     string         `gorm:"size:100" json:"request_id"`
	TransactionID string         `gorm:"size:100" json:"transaction_id"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	RequestData   string         `gorm:"type:text" json:"request_data"`
	ResponseData  string         `gorm:"type:text" json:"response_data"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message"`
	DurationMs    int            `json:"duration_ms"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ProviderConfig struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Provider    string         `gorm:"size:50;not null" json:"provider"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	APIURL      string         `gorm:"size:500;not null" json:"api_url"`
	APIKey      string         `gorm:"size:500" json:"api_key"`
	APISecret   string         `gorm:"size:500" json:"api_secret"`
	RateLimit   int            `gorm:"default:100" json:"rate_limit"`
	TimeoutSec  int            `gorm:"default:30" json:"timeout_sec"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Priority    int            `gorm:"default:1" json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}