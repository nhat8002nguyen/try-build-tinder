package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 50 {
		limit = 50
	}

	notifications, err := h.notificationService.GetUserNotifications(userID, limit, offset)
	if err != nil {
		utils.InternalError(c, "Failed to fetch notifications")
		return
	}

	unreadCount, _ := h.notificationService.GetUnreadCount(userID)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	notifIDParam := c.Param("id")
	notifID, err := uuid.Parse(notifIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := h.notificationService.MarkAsRead(notifID, userID); err != nil {
		utils.InternalError(c, "Failed to mark notification as read")
		return
	}

	utils.MessageResponse(c, http.StatusOK, "Notification marked as read")
}
