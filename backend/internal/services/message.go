package services

import (
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type MessageService struct {
	db *database.Database
}

func NewMessageService(db *database.Database) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) CreateMessage(matchID, senderID uuid.UUID, content string) (*models.Message, error) {
	message := &models.Message{
		MatchID:  matchID,
		SenderID: senderID,
		Content:  content,
	}

	if err := s.db.DB.Create(message).Error; err != nil {
		return nil, err
	}

	if err := s.db.DB.Preload("Sender").First(message, message.ID).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetMessages(matchID uuid.UUID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	query := s.db.DB.Where("match_id = ?", matchID).
		Preload("Sender").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (s *MessageService) GetLatestMessages(matchID uuid.UUID, limit int) ([]models.Message, error) {
	return s.GetMessages(matchID, limit, 0)
}

func (s *MessageService) MarkAsRead(matchID, userID uuid.UUID) error {
	return s.db.DB.Model(&models.Message{}).
		Where("match_id = ? AND sender_id != ? AND is_read = false", matchID, userID).
		Update("is_read", true).Error
}

func (s *MessageService) GetUnreadCount(matchID, userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.DB.Model(&models.Message{}).
		Where("match_id = ? AND sender_id != ? AND is_read = false", matchID, userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MessageService) GetTotalUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.DB.Model(&models.Message{}).
		Joins("JOIN matches ON messages.match_id = matches.id").
		Where("(matches.user1_id = ? OR matches.user2_id = ?) AND messages.sender_id != ? AND messages.is_read = false",
			userID, userID, userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
