package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/websocket"
)

type WebSocketHandler struct {
	hub         *websocket.Hub
	authService *services.AuthService
}

func NewWebSocketHandler(hub *websocket.Hub, authService *services.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	websocket.ServeWs(h.hub, c.Writer, c.Request, userID)
}
