package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/config"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	cfg            *config.Config
	db             *database.Database
	googleConfig   *oauth2.Config
	facebookConfig *oauth2.Config
}

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewAuthService(cfg *config.Config, db *database.Database) *AuthService {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	facebookConfig := &oauth2.Config{
		ClientID:     cfg.FacebookClientID,
		ClientSecret: cfg.FacebookClientSecret,
		RedirectURL:  cfg.FacebookRedirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	return &AuthService{
		cfg:            cfg,
		db:             db,
		googleConfig:   googleConfig,
		facebookConfig: facebookConfig,
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateTokens(user *models.User) (*TokenPair, error) {
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTExpireHours) * time.Hour)

	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.cfg.JWTExpireHours * 3600),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) GetOAuthURL(provider string) (string, error) {
	switch provider {
	case "google":
		return s.googleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline), nil
	case "facebook":
		return s.facebookConfig.AuthCodeURL("state"), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *AuthService) HandleOAuthCallback(ctx context.Context, provider, code string) (*models.User, *TokenPair, error) {
	var oauthConfig *oauth2.Config
	var userInfoURL string

	switch provider {
	case "google":
		oauthConfig = s.googleConfig
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "facebook":
		oauthConfig = s.facebookConfig
		userInfoURL = "https://graph.facebook.com/me?fields=id,name,email"
	default:
		return nil, nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	var oauthAccount models.OAuthAccount
	result := s.db.DB.Where("provider = ? AND provider_user_id = ?", provider, userInfo.ID).First(&oauthAccount)

	var user *models.User

	if result.Error == nil {
		user = &models.User{}
		s.db.DB.First(user, oauthAccount.UserID)

		oauthAccount.AccessToken = token.AccessToken
		oauthAccount.RefreshToken = token.RefreshToken
		s.db.DB.Save(&oauthAccount)
	} else {
		var existingUser models.User
		if err := s.db.DB.Where("email = ?", userInfo.Email).First(&existingUser).Error; err == nil {
			user = &existingUser

			oauthAccount = models.OAuthAccount{
				UserID:         user.ID,
				Provider:       provider,
				ProviderUserID: userInfo.ID,
				AccessToken:    token.AccessToken,
				RefreshToken:   token.RefreshToken,
			}
			s.db.DB.Create(&oauthAccount)
		} else {
			birthdate := time.Now().AddDate(-20, 0, 0)
			user = &models.User{
				Email:      userInfo.Email,
				Name:       userInfo.Name,
				Gender:     models.GenderMale,
				Birthdate:  &birthdate,
				IsVerified: true,
				Preferences: models.UserPreferences{
					MinAge:           18,
					MaxAge:           50,
					MaxDistance:      100,
					GenderPreference: []models.Gender{models.GenderMale, models.GenderFemale},
				},
			}
			if err := s.db.DB.Create(user).Error; err != nil {
				return nil, nil, fmt.Errorf("failed to create user: %w", err)
			}

			oauthAccount = models.OAuthAccount{
				UserID:         user.ID,
				Provider:       provider,
				ProviderUserID: userInfo.ID,
				AccessToken:    token.AccessToken,
				RefreshToken:   token.RefreshToken,
			}
			s.db.DB.Create(&oauthAccount)
		}
	}

	tokens, err := s.GenerateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Register(email, password, name string) (*models.User, *TokenPair, error) {
	var existingUser models.User
	if err := s.db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, nil, errors.New("email already registered")
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	birthdate := time.Now().AddDate(-20, 0, 0)
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
		Gender:       models.GenderMale,
		Birthdate:    &birthdate,
		Preferences: models.UserPreferences{
			MinAge:           18,
			MaxAge:           50,
			MaxDistance:      100,
			GenderPreference: []models.Gender{models.GenderMale, models.GenderFemale},
		},
	}

	if err := s.db.DB.Create(user).Error; err != nil {
		return nil, nil, err
	}

	tokens, err := s.GenerateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(email, password string) (*models.User, *TokenPair, error) {
	var user models.User
	if err := s.db.DB.Preload("Photos").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if user.PasswordHash == "" {
		return nil, nil, errors.New("please login with OAuth provider")
	}

	if !s.CheckPassword(password, user.PasswordHash) {
		return nil, nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastActiveAt = &now
	s.db.DB.Save(&user)

	tokens, err := s.GenerateTokens(&user)
	if err != nil {
		return nil, nil, err
	}

	return &user, tokens, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.User
	if err := s.db.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return s.GenerateTokens(&user)
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.DB.Preload("Photos").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type OAuthUserInfo struct {
	ID    string
	Email string
	Name  string
}

func (s *AuthService) GetGoogleUserInfo(ctx context.Context, code string) (*OAuthUserInfo, error) {
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
