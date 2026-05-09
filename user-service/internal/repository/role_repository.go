package repository

import (
	"github.com/yontech/ppob/user-service/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	return r.FindByID(id)
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("status = ?", "active").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *RoleRepository) AssignRole(userID, roleID uint) error {
	userRole := models.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Create(&userRole).Error
}

func (r *RoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) RemoveUserRole(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}