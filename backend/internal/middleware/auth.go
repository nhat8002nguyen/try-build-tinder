package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tinder-clone/backend/internal/cookies"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
)

func AuthRequired(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(cookies.AccessTokenCookie)
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					token = parts[1]
				}
			}
		}
		if token == "" {
			utils.Unauthorized(c, "Authorization required")
			c.Abort()
			return
		}

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
