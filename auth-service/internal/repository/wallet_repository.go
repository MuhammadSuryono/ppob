package repository

import (
	"github.com/yontech/ppob/auth-service/internal/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *WalletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}