package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// ListNotifications returns notifications for current user
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	unreadOnly := c.DefaultQuery("unread_only", "false") == "true"

	notifs, _, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID.(uint), limit, offset, unreadOnly)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifs})
}

// GetUnreadCount returns count of unread notifications
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID.(uint))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}
	c.JSON(http.StatusOK, dto.UnreadCountResponse{UnreadCount: int(count)})
}

// MarkNotificationRead marks a single notification as read
func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}
	if err := h.notificationService.MarkAsRead(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllNotificationsRead marks all notifications as read for current user
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID.(uint)); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}
