package services

import (
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type SwipeService struct {
	db *database.Database
}

func NewSwipeService(db *database.Database) *SwipeService {
	return &SwipeService{db: db}
}

type SwipeResult struct {
	Swipe   *models.Swipe
	IsMatch bool
	Match   *models.Match
}

func (s *SwipeService) CreateSwipe(swiperID, targetID uuid.UUID, direction models.SwipeDirection) (*SwipeResult, error) {
	var existingSwipe models.Swipe
	if err := s.db.DB.Where("swiper_id = ? AND target_id = ?", swiperID, targetID).First(&existingSwipe).Error; err == nil {
		existingSwipe.Direction = direction
		s.db.DB.Save(&existingSwipe)
		return &SwipeResult{Swipe: &existingSwipe, IsMatch: false}, nil
	}

	swipe := &models.Swipe{
		SwiperID:  swiperID,
		TargetID:  targetID,
		Direction: direction,
	}

	if err := s.db.DB.Create(swipe).Error; err != nil {
		return nil, err
	}

	result := &SwipeResult{Swipe: swipe, IsMatch: false}

	if direction == models.SwipeLike {
		var reverseSwipe models.Swipe
		if err := s.db.DB.Where("swiper_id = ? AND target_id = ? AND direction = ?",
			targetID, swiperID, models.SwipeLike).First(&reverseSwipe).Error; err == nil {
			result.IsMatch = true
		}
	}

	return result, nil
}

func (s *SwipeService) HasSwiped(swiperID, targetID uuid.UUID) (bool, error) {
	var count int64
	if err := s.db.DB.Model(&models.Swipe{}).
		Where("swiper_id = ? AND target_id = ?", swiperID, targetID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SwipeService) GetSwipe(swiperID, targetID uuid.UUID) (*models.Swipe, error) {
	var swipe models.Swipe
	if err := s.db.DB.Where("swiper_id = ? AND target_id = ?", swiperID, targetID).First(&swipe).Error; err != nil {
		return nil, err
	}
	return &swipe, nil
}

func (s *SwipeService) GetLikesReceived(userID uuid.UUID) ([]models.Swipe, error) {
	var swipes []models.Swipe
	if err := s.db.DB.Where("target_id = ? AND direction = ?", userID, models.SwipeLike).
		Preload("Swiper").
		Order("created_at DESC").
		Find(&swipes).Error; err != nil {
		return nil, err
	}
	return swipes, nil
}
