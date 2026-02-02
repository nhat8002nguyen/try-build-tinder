package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/models"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
	"github.com/tinder-clone/backend/internal/websocket"
)

type SwipeHandler struct {
	swipeService        *services.SwipeService
	matchService        *services.MatchService
	notificationService *services.NotificationService
	wsHub               *websocket.Hub
}

func NewSwipeHandler(
	swipeService *services.SwipeService,
	matchService *services.MatchService,
	notificationService *services.NotificationService,
	wsHub *websocket.Hub,
) *SwipeHandler {
	return &SwipeHandler{
		swipeService:        swipeService,
		matchService:        matchService,
		notificationService: notificationService,
		wsHub:               wsHub,
	}
}

type CreateSwipeRequest struct {
	TargetID  string `json:"target_id" binding:"required"`
	Direction string `json:"direction" binding:"required,oneof=like dislike"`
}

type SwipeResponse struct {
	IsMatch bool          `json:"is_match"`
	Match   *models.Match `json:"match,omitempty"`
}

func (h *SwipeHandler) CreateSwipe(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var req CreateSwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		utils.BadRequest(c, "Invalid target ID")
		return
	}

	if userID == targetID {
		utils.BadRequest(c, "Cannot swipe on yourself")
		return
	}

	direction := models.SwipeDirection(req.Direction)
	result, err := h.swipeService.CreateSwipe(userID, targetID, direction)
	if err != nil {
		utils.InternalError(c, "Failed to record swipe")
		return
	}

	response := SwipeResponse{IsMatch: result.IsMatch}

	if result.IsMatch {
		match, err := h.matchService.CreateMatch(userID, targetID)
		if err != nil {
			utils.InternalError(c, "Failed to create match")
			return
		}
		response.Match = match

		h.notifyMatch(userID, targetID, match)
	}

	utils.SuccessResponse(c, http.StatusCreated, response)
}

func (h *SwipeHandler) notifyMatch(user1ID, user2ID uuid.UUID, match *models.Match) {
	h.wsHub.SendToUser(user1ID, websocket.Message{
		Type: "match",
		Payload: map[string]interface{}{
			"match_id":      match.ID.String(),
			"other_user_id": user2ID.String(),
		},
	})

	h.wsHub.SendToUser(user2ID, websocket.Message{
		Type: "match",
		Payload: map[string]interface{}{
			"match_id":      match.ID.String(),
			"other_user_id": user1ID.String(),
		},
	})
}
