package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
	"github.com/tinder-clone/backend/internal/websocket"
)

type MessageHandler struct {
	messageService *services.MessageService
	matchService   *services.MatchService
	wsHub          *websocket.Hub
}

func NewMessageHandler(
	messageService *services.MessageService,
	matchService *services.MatchService,
	wsHub *websocket.Hub,
) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		matchService:   matchService,
		wsHub:          wsHub,
	}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	matchIDParam := c.Param("id")
	matchID, err := uuid.Parse(matchIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid match ID")
		return
	}

	isUserInMatch, err := h.matchService.IsUserInMatch(userID, matchID)
	if err != nil || !isUserInMatch {
		utils.Forbidden(c, "Not authorized to view these messages")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	messages, err := h.messageService.GetMessages(matchID, limit, offset)
	if err != nil {
		utils.InternalError(c, "Failed to fetch messages")
		return
	}

	h.messageService.MarkAsRead(matchID, userID)

	utils.SuccessResponse(c, http.StatusOK, messages)
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	matchIDParam := c.Param("id")
	matchID, err := uuid.Parse(matchIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid match ID")
		return
	}

	isUserInMatch, err := h.matchService.IsUserInMatch(userID, matchID)
	if err != nil || !isUserInMatch {
		utils.Forbidden(c, "Not authorized to send messages to this match")
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	message, err := h.messageService.CreateMessage(matchID, userID, req.Content)
	if err != nil {
		utils.InternalError(c, "Failed to send message")
		return
	}

	h.matchService.UpdateLastMessage(matchID)

	match, _ := h.matchService.GetMatch(matchID)
	if match != nil {
		recipientID := match.GetOtherUserID(userID)
		h.wsHub.SendToUser(recipientID, websocket.Message{
			Type: "message",
			Payload: map[string]interface{}{
				"message":  message,
				"match_id": matchID.String(),
			},
		})
	}

	utils.SuccessResponse(c, http.StatusCreated, message)
}
