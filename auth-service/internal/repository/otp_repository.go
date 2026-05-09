package repository

import (
	"github.com/yontech/ppob/auth-service/internal/models"
	"gorm.io/gorm"
	"time"
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(otp *models.OTP) error {
	return r.db.Create(otp).Error
}

func (r *OTPRepository) FindValidOTP(userID uint, code, otpType string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.Where("user_id = ? AND code = ? AND type = ? AND used_at IS NULL AND expires_at > ?",
		userID, code, otpType, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) MarkAsUsed(otp *models.OTP) error {
	now := time.Now()
	otp.UsedAt = &now
	return r.db.Save(otp).Error
}

func (r *OTPRepository) DeleteExpired(userID uint, otpType string) error {
	return r.db.Where("user_id = ? AND type = ? AND expires_at < ?", userID, otpType, time.Now()).Delete(&models.OTP{}).Error
}

func (r *OTPRepository) FindLatestOTP(userID uint, otpType string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.Where("user_id = ? AND type = ?", userID, otpType).Order("created_at DESC").First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}