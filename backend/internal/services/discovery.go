package services

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
)

type DiscoveryService struct {
	db *database.Database
}

func NewDiscoveryService(db *database.Database) *DiscoveryService {
	return &DiscoveryService{db: db}
}

type DiscoveryFilters struct {
	MinAge           int
	MaxAge           int
	MaxDistance      int
	GenderPreference []models.Gender
	Limit            int
}

type ScoredUser struct {
	User  models.User
	Score float64
}

func (s *DiscoveryService) GetPotentialMatches(userID uuid.UUID, filters DiscoveryFilters) ([]models.User, error) {
	var currentUser models.User
	if err := s.db.DB.First(&currentUser, userID).Error; err != nil {
		return nil, err
	}

	if filters.MinAge == 0 {
		filters.MinAge = currentUser.Preferences.MinAge
	}
	if filters.MaxAge == 0 {
		filters.MaxAge = currentUser.Preferences.MaxAge
	}
	if filters.MaxDistance == 0 {
		filters.MaxDistance = currentUser.Preferences.MaxDistance
	}
	if len(filters.GenderPreference) == 0 {
		filters.GenderPreference = currentUser.Preferences.GenderPreference
	}
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	now := time.Now()
	minBirthdate := now.AddDate(-filters.MaxAge-1, 0, 0)
	maxBirthdate := now.AddDate(-filters.MinAge, 0, 0)

	var swipedUserIDs []uuid.UUID
	s.db.DB.Model(&models.Swipe{}).
		Where("swiper_id = ?", userID).
		Pluck("target_id", &swipedUserIDs)

	var matchedUserIDs []uuid.UUID
	var matches []models.Match
	s.db.DB.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&matches)
	for _, m := range matches {
		if m.User1ID == userID {
			matchedUserIDs = append(matchedUserIDs, m.User2ID)
		} else {
			matchedUserIDs = append(matchedUserIDs, m.User1ID)
		}
	}

	excludeIDs := append(swipedUserIDs, matchedUserIDs...)
	excludeIDs = append(excludeIDs, userID)

	var candidates []models.User
	query := s.db.DB.Where("id NOT IN ?", excludeIDs).
		Where("is_active = ?", true).
		Where("birthdate IS NOT NULL").
		Where("birthdate BETWEEN ? AND ?", minBirthdate, maxBirthdate)

	if len(filters.GenderPreference) > 0 {
		query = query.Where("gender IN ?", filters.GenderPreference)
	}

	query = query.Preload("Photos").
		Limit(filters.Limit * 3)

	if err := query.Find(&candidates).Error; err != nil {
		return nil, err
	}

	scoredUsers := make([]ScoredUser, 0, len(candidates))
	for _, candidate := range candidates {
		distance := calculateDistance(
			currentUser.Location.Latitude, currentUser.Location.Longitude,
			candidate.Location.Latitude, candidate.Location.Longitude,
		)

		if filters.MaxDistance > 0 && distance > float64(filters.MaxDistance) {
			continue
		}

		score := s.calculateScore(currentUser, candidate, distance)
		scoredUsers = append(scoredUsers, ScoredUser{User: candidate, Score: score})
	}

	sort.Slice(scoredUsers, func(i, j int) bool {
		return scoredUsers[i].Score > scoredUsers[j].Score
	})

	result := make([]models.User, 0, filters.Limit)
	for i := 0; i < len(scoredUsers) && i < filters.Limit; i++ {
		result = append(result, scoredUsers[i].User)
	}

	return result, nil
}

func (s *DiscoveryService) calculateScore(currentUser models.User, candidate models.User, distance float64) float64 {
	activityScore := s.calculateActivityScore(candidate)
	distanceScore := s.calculateDistanceScore(distance)
	profileScore := s.calculateProfileScore(candidate)

	w1, w2, w3 := 0.4, 0.35, 0.25
	return w1*activityScore + w2*distanceScore + w3*profileScore
}

func (s *DiscoveryService) calculateActivityScore(user models.User) float64 {
	if user.LastActiveAt == nil {
		return 0.1
	}
	
	hoursSinceActive := time.Since(*user.LastActiveAt).Hours()

	switch {
	case hoursSinceActive < 1:
		return 1.0
	case hoursSinceActive < 24:
		return 0.8
	case hoursSinceActive < 72:
		return 0.6
	case hoursSinceActive < 168:
		return 0.4
	default:
		return 0.2
	}
}

func (s *DiscoveryService) calculateDistanceScore(distance float64) float64 {
	if distance <= 5 {
		return 1.0
	} else if distance <= 25 {
		return 0.8
	} else if distance <= 50 {
		return 0.6
	} else if distance <= 100 {
		return 0.4
	}
	return 0.2
}

func (s *DiscoveryService) calculateProfileScore(user models.User) float64 {
	score := 0.0

	if len(user.Photos) > 0 {
		photoScore := float64(len(user.Photos)) * 0.15
		if photoScore > 0.45 {
			photoScore = 0.45
		}
		score += photoScore
	}

	if user.Bio != "" {
		bioLength := len(user.Bio)
		if bioLength > 50 {
			score += 0.3
		} else if bioLength > 20 {
			score += 0.2
		} else {
			score += 0.1
		}
	}

	if user.IsVerified {
		score += 0.25
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
