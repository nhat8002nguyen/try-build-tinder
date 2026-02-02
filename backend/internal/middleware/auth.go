package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
)

func AuthRequired(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}
