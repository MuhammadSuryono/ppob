package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Phone        string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	FullName     string         `gorm:"size:255;not null" json:"full_name"`
	PIN          string         `gorm:"size:255" json:"-"`
	Role         string         `gorm:"size:20;default:user" json:"role"`
	Status       string         `gorm:"size:20;default:active" json:"status"`
	EmailVerified bool          `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool          `gorm:"default:false" json:"phone_verified"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	Type      string    `gorm:"size:20;not null" json:"type"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
}

type DeviceFingerprint struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	FingerprintHash string    `gorm:"size:64;not null" json:"fingerprint_hash"`
	UserAgent       string    `gorm:"type:text" json:"user_agent"`
	IPAddress       string    `gorm:"size:50" json:"ip_address"`
	TrustScore      int       `gorm:"default:0" json:"trust_score"`
	IsTrusted       bool      `gorm:"default:false" json:"is_trusted"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
}

type Wallet struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance         float64        `gorm:"default:0" json:"balance"`
	HoldAmount      float64        `gorm:"default:0" json:"hold_amount"`
	Currency        string         `gorm:"size:10;default:IDR" json:"currency"`
	Status          string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;size:500;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}