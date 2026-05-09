package services

import (
	"context"
	"errors"
	"time"

	"github.com/yontech/ppob/transaction-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDailyLimitExceeded = errors.New("daily limit exceeded")
)

type DailyLimitService struct {
	db *gorm.DB
}

func NewDailyLimitService(db *gorm.DB) *DailyLimitService {
	return &DailyLimitService{db: db}
}

func (s *DailyLimitService) CheckAndUpdateLimit(ctx context.Context, userID uint, amount float64) error {
	today := time.Now().Format("2006-01-02")

	var limit models.DailyLimit
	err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&limit).Error

	if err == gorm.ErrRecordNotFound {
		limit = models.DailyLimit{
			UserID:    userID,
			Date:      today,
			Count:     0,
			TotalAmount: 0,
		}
		if err := s.db.Create(&limit).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if limit.MaxCount > 0 && limit.Count >= limit.MaxCount {
		return ErrDailyLimitExceeded
	}

	if limit.MaxAmount > 0 && limit.TotalAmount+amount > limit.MaxAmount {
		return ErrDailyLimitExceeded
	}

	return s.db.Model(&limit).Updates(map[string]interface{}{
		"count":        gorm.Expr("count + 1"),
		"total_amount": gorm.Expr("total_amount + ?", amount),
	}).Error
}

func (s *DailyLimitService) GetTodayLimit(userID uint) (*models.DailyLimit, error) {
	today := time.Now().Format("2006-01-02")

	var limit models.DailyLimit
	err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &limit, nil
}

func (s *DailyLimitService) SetDailyLimit(userID uint, maxCount int, maxAmount float64) error {
	today := time.Now().Format("2006-01-02")

	var limit models.DailyLimit
	err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&limit).Error

	if err == gorm.ErrRecordNotFound {
		limit = models.DailyLimit{
			UserID:     userID,
			Date:       today,
			MaxCount:   maxCount,
			MaxAmount:  maxAmount,
			Count:      0,
			TotalAmount: 0,
		}
		return s.db.Create(&limit).Error
	}

	limit.MaxCount = maxCount
	limit.MaxAmount = maxAmount
	return s.db.Save(&limit).Error
}

func (s *DailyLimitService) GetRemainingLimit(userID uint) (remainingCount int, remainingAmount float64, err error) {
	limit, err := s.GetTodayLimit(userID)
	if err != nil {
		return 0, 0, err
	}

	if limit == nil {
		return -1, -1, nil
	}

	remainingCount = limit.MaxCount - limit.Count
	if remainingCount < 0 {
		remainingCount = 0
	}

	remainingAmount = limit.MaxAmount - limit.TotalAmount
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	return remainingCount, remainingAmount, nil
}

func (s *DailyLimitService) GetDailyLimitSummary(userID uint, date string) (map[string]interface{}, error) {
	var limit models.DailyLimit
	err := s.db.Where("user_id = ? AND date = ?", userID, date).First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"user_id":    userID,
				"date":       date,
				"count":      0,
				"total_amount": 0.0,
				"max_count":  0,
				"max_amount": 0.0,
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"user_id":       limit.UserID,
		"date":          limit.Date,
		"count":         limit.Count,
		"total_amount":  limit.TotalAmount,
		"max_count":     limit.MaxCount,
		"max_amount":    limit.MaxAmount,
		"remaining_count": limit.MaxCount - limit.Count,
		"remaining_amount": limit.MaxAmount - limit.TotalAmount,
	}, nil
}