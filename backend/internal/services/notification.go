package services

import (
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type NotificationService struct {
	db *database.Database
}

func NewNotificationService(db *database.Database) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(userID uuid.UUID, notifType models.NotificationType, payload models.NotificationPayload) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Payload: payload,
	}

	if err := s.db.DB.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) CreateMatchNotification(userID, matchID, otherUserID uuid.UUID, otherUserName string) (*models.Notification, error) {
	payload := models.NotificationPayload{
		"match_id":        matchID.String(),
		"other_user_id":   otherUserID.String(),
		"other_user_name": otherUserName,
		"message":         "You have a new match with " + otherUserName + "!",
	}
	return s.Create(userID, models.NotificationTypeMatch, payload)
}

func (s *NotificationService) CreateMessageNotification(userID, matchID, senderID uuid.UUID, senderName, messagePreview string) (*models.Notification, error) {
	if len(messagePreview) > 50 {
		messagePreview = messagePreview[:47] + "..."
	}
	payload := models.NotificationPayload{
		"match_id":    matchID.String(),
		"sender_id":   senderID.String(),
		"sender_name": senderName,
		"preview":     messagePreview,
		"message":     senderName + ": " + messagePreview,
	}
	return s.Create(userID, models.NotificationTypeMessage, payload)
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := s.db.DB.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *NotificationService) GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := s.db.DB.Where("user_id = ? AND is_read = false", userID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) MarkAsRead(notificationID, userID uuid.UUID) error {
	return s.db.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.db.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *NotificationService) Delete(notificationID, userID uuid.UUID) error {
	return s.db.DB.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error
}
