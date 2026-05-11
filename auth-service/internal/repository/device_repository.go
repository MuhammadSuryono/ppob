package repository

import (
	"github.com/yontech/ppob/auth-service/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) FindByFingerprint(userID uint, fingerprint string) (*models.DeviceFingerprint, error) {
	var device models.DeviceFingerprint
	err := r.db.Where("user_id = ? AND fingerprint_hash = ?", userID, fingerprint).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Create(device *models.DeviceFingerprint) error {
	return r.db.Create(device).Error
}

func (r *DeviceRepository) Update(device *models.DeviceFingerprint) error {
	return r.db.Save(device).Error
}
