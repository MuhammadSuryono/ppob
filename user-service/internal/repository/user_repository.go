package repository

import (
	"errors"
	"time"

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

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
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

func (r *UserRepository) GetUserRoles(userID uint) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.db.Where("user_id = ?", userID).Find(&roles).Error
	return roles, err
}

type WalletModel struct {
	ID       uint    `gorm:"primaryKey"`
	UserID   uint    `gorm:"index"`
	Balance  float64 `gorm:"default:0"`
}

func (r *UserRepository) GetWalletByUserAndRole(userID, roleID uint) (*WalletModel, error) {
	var wallet WalletModel
	err := r.db.Table("wallets").
		Where("user_id = ?", userID).
		Select("id, user_id, balance").
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *UserRepository) GetWalletBalance(userID uint) (float64, error) {
	var wallet WalletModel
	err := r.db.Table("wallets").Where("user_id = ?", userID).Select("balance").First(&wallet).Error
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (r *UserRepository) FindStaffUsers(mitraID uint, limit, offset int, status string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN staff_global_margin_settings ON users.id = staff_global_margin_settings.staff_id").
		Where("roles.name = ?", "staff")

	if mitraID > 0 {
		query = query.Where("staff_global_margin_settings.mitra_id = ?", mitraID)
	}

	if status != "" && status != "all" {
		query = query.Where("users.status = ?", status)
	}

	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Order("users.created_at DESC").Find(&users).Error

	return users, total, err
}

func (r *UserRepository) FindStaffByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Model(&models.User{}).Joins("JOIN user_roles ON users.id = user_roles.user_id").Joins("JOIN roles ON user_roles.role_id = roles.id").Where("users.id = ? AND roles.name = ?", id, "staff").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type staffTodayStats struct {
	Count  int64   `gorm:"column:count"`
	Amount float64 `gorm:"column:amount"`
}

func (r *UserRepository) GetStaffTodayStats(staffID uint) (int, float64, error) {
	var stats staffTodayStats
	dateStr := time.Now().Format("2006-01-02")

	err := r.db.Table("transactions").
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("user_id = ? AND status = ? AND DATE(created_at) = ?", staffID, "success", dateStr).
		Scan(&stats).Error
	if err != nil {
		return 0, 0, err
	}

	return int(stats.Count), stats.Amount, nil
}

type dailyLimitResult struct {
	MaxCount  int     `json:"max_count"`
	MaxAmount float64 `json:"max_amount"`
}

func (r *UserRepository) GetStaffDailyLimit(staffID uint) (int, float64, error) {
	var limit dailyLimitResult
	dateStr := time.Now().Format("2006-01-02")

	err := r.db.Table("daily_limits").
		Where("user_id = ? AND date = ?", staffID, dateStr).
		Select("max_count, max_amount").
		Scan(&limit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	return limit.MaxCount, limit.MaxAmount, nil
}
