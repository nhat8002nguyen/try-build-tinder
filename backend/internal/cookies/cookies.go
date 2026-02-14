package cookies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tinder-clone/backend/internal/config"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

func SetAuthCookies(c *gin.Context, accessToken, refreshToken string, cfg *config.Config) {
	secure := cfg.Environment == "production"
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	accessMaxAge := cfg.JWTExpireHours * 3600
	refreshMaxAge := 7 * 24 * 3600

	accessCookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessMaxAge,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   refreshMaxAge,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)
}

func ClearAuthCookies(c *gin.Context) {
	expired := time.Unix(0, 0)
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  expired,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  expired,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)
}
