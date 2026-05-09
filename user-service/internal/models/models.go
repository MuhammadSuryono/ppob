package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Email         string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Phone         string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	FullName      string         `gorm:"size:255;not null" json:"full_name"`
	Role          string         `gorm:"size:20;default:user" json:"role"`
	Status        string         `gorm:"size:20;default:active" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool           `gorm:"default:false" json:"phone_verified"`
	Avatar        string         `gorm:"size:500" json:"avatar"`
	Address       string         `gorm:"type:text" json:"address"`
	DateOfBirth   *time.Time     `json:"date_of_birth"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Permissions string         `gorm:"type:text" json:"permissions"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRole struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	RoleID     uint      `gorm:"index;not null" json:"role_id"`
	AssignedBy *uint     `gorm:"index" json:"assigned_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StaffGlobalMarginSetting struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	MitraID         uint           `gorm:"index;not null" json:"mitra_id"`
	StaffID         uint           `gorm:"index;not null" json:"staff_id"`
	SchemeType      string         `gorm:"size:20;not null;default:FixedAllowance" json:"scheme_type"`
	GlobalMarginPercent float64    `gorm:"type:decimal(5,2);default:0" json:"global_margin_percent"`
	FixedAllowance float64        `gorm:"type:decimal(15,2);default:0" json:"fixed_allowance"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type StaffProductMarginOverride struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	MitraID         uint           `gorm:"index;not null" json:"mitra_id"`
	StaffID         uint           `gorm:"index;not null" json:"staff_id"`
	ProductID       uint           `gorm:"index;not null" json:"product_id"`
	ProductCode     string         `gorm:"size:100" json:"product_code"`
	MarginPercent   float64        `gorm:"type:decimal(5,2);default:0" json:"margin_percent"`
	FixedMargin     float64        `gorm:"type:decimal(15,2);default:0" json:"fixed_margin"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}