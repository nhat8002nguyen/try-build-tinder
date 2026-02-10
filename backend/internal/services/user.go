package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/apperrors"
	"github.com/tinder-clone/backend/internal/constants"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type UserService struct {
	db *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.DB.Preload("Photos").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByIDSimple(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.DB.Preload("Photos").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type UpdateProfileInput struct {
	Name        *string                  `json:"name"`
	Gender      *models.Gender           `json:"gender"`
	Birthdate   *time.Time               `json:"birthdate"`
	Bio         *string                  `json:"bio"`
	Preferences *models.UserPreferences  `json:"preferences"`
}

func (s *UserService) UpdateProfile(userID uuid.UUID, input UpdateProfileInput) (*models.User, error) {
	var user models.User
	if err := s.db.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Gender != nil {
		user.Gender = *input.Gender
	}
	if input.Birthdate != nil {
		user.Birthdate = input.Birthdate
	}
	if input.Bio != nil {
		user.Bio = *input.Bio
	}
	if input.Preferences != nil {
		user.Preferences = *input.Preferences
	}

	user.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetByIDSimple(userID)
}

func (s *UserService) UpdateLocation(userID uuid.UUID, location models.Location) error {
	return s.db.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"location":       location,
			"last_active_at": time.Now(),
		}).Error
}

func (s *UserService) UpdateLastActive(userID uuid.UUID) error {
	return s.db.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_active_at", time.Now()).Error
}

func (s *UserService) AddPhoto(userID uuid.UUID, photoURL string, order int) (*models.UserPhoto, error) {
	count, err := s.GetPhotoCount(userID)
	if err != nil {
		return nil, err
	}
	if count >= int64(constants.MaxPhotosPerUser) {
		return nil, apperrors.ErrMaxPhotosReached
	}

	photo := &models.UserPhoto{
		UserID:       userID,
		PhotoURL:     photoURL,
		DisplayOrder: order,
		IsApproved:   true,
	}

	if err := s.db.DB.Create(photo).Error; err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *UserService) DeletePhoto(userID, photoID uuid.UUID) error {
	return s.db.DB.Where("id = ? AND user_id = ?", photoID, userID).Delete(&models.UserPhoto{}).Error
}

func (s *UserService) GetUserPhotos(userID uuid.UUID) ([]models.UserPhoto, error) {
	var photos []models.UserPhoto
	if err := s.db.DB.Where("user_id = ?", userID).Order("display_order ASC").Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *UserService) GetPhotoCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.DB.Model(&models.UserPhoto{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
