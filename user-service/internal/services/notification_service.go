package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/models"
	"github.com/yontech/ppob/user-service/internal/repository"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uint, notifType, title, message string) error {
	notif := &models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return s.notificationRepo.Create(notif)
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uint, limit, offset int, unreadOnly bool) ([]dto.NotificationResponse, int64, error) {
	notifs, total, err := s.notificationRepo.GetByUser(userID, limit, offset, unreadOnly)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		responses[i] = dto.NotificationResponse{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		}
	}
	return responses, total, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.notificationRepo.CountUnread(userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notifID uint, userID uint) error {
	return s.notificationRepo.MarkRead(notifID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.notificationRepo.MarkAllRead(userID)
}

// SendStaffInvitationNotification creates a notification for a Mitra about staff invitation
func (s *NotificationService) SendStaffInvitationNotification(ctx context.Context, mitraID uint, staffName string) error {
	title := "Undangan Staff"
	message := fmt.Sprintf("Staff %s telah berhasil ditambahkan.", staffName)
	return s.CreateNotification(ctx, mitraID, "staff_invite", title, message)
}
