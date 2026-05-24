package repository

import (
	"github.com/yontech/ppob/auth-service/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) AssignRole(userID uint, roleID uint) error {
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(userRole).Error
}
