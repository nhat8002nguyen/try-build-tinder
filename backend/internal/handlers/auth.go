package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.BadRequest(c, "Invalid email format")
		return
	}

	if !utils.IsValidPassword(req.Password) {
		utils.BadRequest(c, "Password must be at least 8 characters")
		return
	}

	user, tokens, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tokens, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, tokens)
}

func (h *AuthHandler) OAuthRedirect(c *gin.Context) {
	provider := c.Param("provider")

	url, err := h.authService.GetOAuthURL(provider)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	if code == "" {
		utils.BadRequest(c, "Authorization code required")
		return
	}

	_, tokens, err := h.authService.HandleOAuthCallback(c.Request.Context(), provider, code)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	frontendURL := "http://localhost:3000/auth/callback"
	redirectURL := frontendURL + "?access_token=" + tokens.AccessToken + "&refresh_token=" + tokens.RefreshToken

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}
