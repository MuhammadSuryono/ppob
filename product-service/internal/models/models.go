package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Code          string         `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Brand         string         `gorm:"size:100" json:"brand"`
	CategoryID    uint           `gorm:"index;not null" json:"category_id"`
	Provider      string         `gorm:"size:50" json:"provider"`
	Price         float64        `gorm:"not null" json:"price"`
	PriceAPI      float64        `json:"price_api"`
	Stock         int            `gorm:"default:-1" json:"stock"`
	Status        string         `gorm:"size:20;default:active" json:"status"`
	Description   string         `gorm:"type:text" json:"description"`
	LastSyncAt    *time.Time     `json:"last_sync_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Code            string         `gorm:"size:50;not null" json:"code"`
	Description     string         `gorm:"type:text" json:"description"`
	Icon            string         `gorm:"size:255" json:"icon"`
	SortOrder       int            `gorm:"default:0" json:"sort_order"`
	Status          string         `gorm:"size:20;default:active" json:"status"`
	InputType       string         `gorm:"size:50;default:TEXT" json:"input_type"`
	InputLabel      string         `gorm:"size:100" json:"input_label"`
	Placeholder     string         `gorm:"size:255" json:"placeholder"`
	ValidationRegex string         `gorm:"size:255" json:"validation_regex"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}