package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

type MatchHandler struct {
	matchService *services.MatchService
}

func NewMatchHandler(matchService *services.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

type MatchResponse struct {
	ID            string      `json:"id"`
	OtherUser     interface{} `json:"other_user"`
	MatchedAt     string      `json:"matched_at"`
	LastMessageAt *string     `json:"last_message_at"`
}

func (h *MatchHandler) GetMatches(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	matches, err := h.matchService.GetUserMatches(userID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch matches")
		return
	}

	response := make([]MatchResponse, 0, len(matches))
	for _, match := range matches {
		otherUser := match.GetOtherUser(userID)

		var lastMessageAt *string
		if match.LastMessageAt != nil {
			formatted := match.LastMessageAt.Format("2006-01-02T15:04:05Z07:00")
			lastMessageAt = &formatted
		}

		response = append(response, MatchResponse{
			ID:            match.ID.String(),
			OtherUser:     otherUser,
			MatchedAt:     match.MatchedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastMessageAt: lastMessageAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	matchIDParam := c.Param("id")
	matchID, err := uuid.Parse(matchIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid match ID")
		return
	}

	match, err := h.matchService.GetMatch(matchID)
	if err != nil {
		utils.NotFound(c, "Match not found")
		return
	}

	isUserInMatch, err := h.matchService.IsUserInMatch(userID, matchID)
	if err != nil || !isUserInMatch {
		utils.Forbidden(c, "Not authorized to view this match")
		return
	}

	otherUser := match.GetOtherUser(userID)

	var lastMessageAt *string
	if match.LastMessageAt != nil {
		formatted := match.LastMessageAt.Format("2006-01-02T15:04:05Z07:00")
		lastMessageAt = &formatted
	}

	response := MatchResponse{
		ID:            match.ID.String(),
		OtherUser:     otherUser,
		MatchedAt:     match.MatchedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastMessageAt: lastMessageAt,
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}
