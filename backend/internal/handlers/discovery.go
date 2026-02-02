package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/models"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

type DiscoveryHandler struct {
	discoveryService *services.DiscoveryService
}

func NewDiscoveryHandler(discoveryService *services.DiscoveryService) *DiscoveryHandler {
	return &DiscoveryHandler{discoveryService: discoveryService}
}

func (h *DiscoveryHandler) GetPotentialMatches(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	filters := services.DiscoveryFilters{}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 50 {
			filters.Limit = limit
		}
	}

	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		if minAge, err := strconv.Atoi(minAgeStr); err == nil && minAge >= 18 {
			filters.MinAge = minAge
		}
	}

	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil && maxAge <= 100 {
			filters.MaxAge = maxAge
		}
	}

	if maxDistStr := c.Query("max_distance"); maxDistStr != "" {
		if maxDist, err := strconv.Atoi(maxDistStr); err == nil && maxDist > 0 {
			filters.MaxDistance = maxDist
		}
	}

	if genderPref := c.Query("gender_preference"); genderPref != "" {
		genders := strings.Split(genderPref, ",")
		for _, g := range genders {
			gender := models.Gender(strings.TrimSpace(g))
			if gender == models.GenderMale || gender == models.GenderFemale || gender == models.GenderOther {
				filters.GenderPreference = append(filters.GenderPreference, gender)
			}
		}
	}

	users, err := h.discoveryService.GetPotentialMatches(userID, filters)
	if err != nil {
		utils.InternalError(c, "Failed to fetch potential matches")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, users)
}
