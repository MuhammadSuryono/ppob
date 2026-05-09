package repository

import (
	"github.com/yontech/ppob/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) List(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

type WalletModel struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index"`
	RoleID uint `gorm:"index"`
}

func (r *UserRepository) GetUserRoles(userID uint) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.db.Where("user_id = ?", userID).Find(&roles).Error
	return roles, err
}

func (r *UserRepository) GetWalletByUserAndRole(userID, roleID uint) (*WalletModel, error) {
	var wallet WalletModel
	err := r.db.Table("wallets").
		Joins("JOIN user_roles ON wallets.user_id = user_roles.user_id AND user_roles.role_id = ?", roleID).
		Where("wallets.user_id = ?", userID).
		Select("wallets.id, wallets.user_id, user_roles.role_id").
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}