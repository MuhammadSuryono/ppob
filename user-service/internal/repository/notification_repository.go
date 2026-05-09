package repository

import (
	"github.com/yontech/ppob/user-service/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notif *models.Notification) error {
	return r.db.Create(notif).Error
}

func (r *NotificationRepository) GetByUser(userID uint, limit, offset int, unreadOnly bool) ([]models.Notification, int64, error) {
	var notifs []models.Notification
	query := r.db.Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	var total int64
	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifs).Error
	return notifs, total, err
}

func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkRead(notifID uint, userID uint) error {
	return r.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", notifID, userID).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}
