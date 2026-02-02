package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type MatchService struct {
	db *database.Database
}

func NewMatchService(db *database.Database) *MatchService {
	return &MatchService{db: db}
}

func (s *MatchService) CreateMatch(user1ID, user2ID uuid.UUID) (*models.Match, error) {
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	var existingMatch models.Match
	if err := s.db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&existingMatch).Error; err == nil {
		return &existingMatch, nil
	}

	match := &models.Match{
		User1ID: user1ID,
		User2ID: user2ID,
	}

	if err := s.db.DB.Create(match).Error; err != nil {
		return nil, err
	}

	return match, nil
}

func (s *MatchService) GetMatch(matchID uuid.UUID) (*models.Match, error) {
	var match models.Match
	if err := s.db.DB.Preload("User1.Photos").Preload("User2.Photos").First(&match, matchID).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *MatchService) GetMatchByUsers(user1ID, user2ID uuid.UUID) (*models.Match, error) {
	var match models.Match
	if err := s.db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *MatchService) GetUserMatches(userID uuid.UUID) ([]models.Match, error) {
	var matches []models.Match
	if err := s.db.DB.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Preload("User1.Photos").
		Preload("User2.Photos").
		Order("COALESCE(last_message_at, matched_at) DESC").
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) IsUserInMatch(userID, matchID uuid.UUID) (bool, error) {
	var match models.Match
	if err := s.db.DB.First(&match, matchID).Error; err != nil {
		return false, err
	}
	return match.User1ID == userID || match.User2ID == userID, nil
}

func (s *MatchService) UpdateLastMessage(matchID uuid.UUID) error {
	now := time.Now()
	return s.db.DB.Model(&models.Match{}).
		Where("id = ?", matchID).
		Update("last_message_at", &now).Error
}

func (s *MatchService) DeleteMatch(userID, matchID uuid.UUID) error {
	var match models.Match
	if err := s.db.DB.First(&match, matchID).Error; err != nil {
		return err
	}

	if match.User1ID != userID && match.User2ID != userID {
		return errors.New("not authorized to delete this match")
	}

	return s.db.DB.Delete(&match).Error
}

func (s *MatchService) GetMatchCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.DB.Model(&models.Match{}).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
