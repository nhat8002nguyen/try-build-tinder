package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/tinder-clone/backend/internal/config"
	"github.com/tinder-clone/backend/internal/database"
	"github.com/tinder-clone/backend/internal/handlers"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/websocket"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	log.Info().Msg("Database migrations completed")

	authService := services.NewAuthService(cfg, db)
	userService := services.NewUserService(db)
	swipeService := services.NewSwipeService(db)
	matchService := services.NewMatchService(db)
	messageService := services.NewMessageService(db)
	notificationService := services.NewNotificationService(db)
	discoveryService := services.NewDiscoveryService(db)
	storageService := services.NewStorageService(cfg)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService, storageService)
	swipeHandler := handlers.NewSwipeHandler(swipeService, matchService, notificationService, wsHub)
	matchHandler := handlers.NewMatchHandler(matchService)
	messageHandler := handlers.NewMessageHandler(messageService, matchService, wsHub)
	discoveryHandler := handlers.NewDiscoveryHandler(discoveryService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	wsHandler := handlers.NewWebSocketHandler(wsHub, authService)

	router := gin.Default()

	corsOrigins := []string{"http://localhost:3000", "http://localhost:5173", "https://spark.vnhatng.com", "http://spark.vnhatng.com"}
	if cfg.CORSAllowedOrigins != "" {
		parts := strings.Split(cfg.CORSAllowedOrigins, ",")
		for _, o := range parts {
			if o = strings.TrimSpace(o); o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}
	originSet := make(map[string]bool)
	for _, o := range corsOrigins {
		originSet[strings.TrimSuffix(o, "/")] = true
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowOriginFunc: func(origin string) bool { return originSet[strings.TrimSuffix(origin, "/")] },
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/oauth/:provider", authHandler.OAuthRedirect)
			auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
			auth.GET("/me", middleware.AuthRequired(authService), authHandler.GetCurrentUser)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthRequired(authService))
		{
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/me", userHandler.UpdateProfile)
			users.POST("/me/photos", userHandler.UploadPhoto)
			users.DELETE("/me/photos/:photoId", userHandler.DeletePhoto)
			users.PUT("/me/location", userHandler.UpdateLocation)
		}

		discover := api.Group("/discover")
		discover.Use(middleware.AuthRequired(authService))
		{
			discover.GET("", discoveryHandler.GetPotentialMatches)
		}

		swipes := api.Group("/swipes")
		swipes.Use(middleware.AuthRequired(authService))
		{
			swipes.POST("", swipeHandler.CreateSwipe)
		}

		matches := api.Group("/matches")
		matches.Use(middleware.AuthRequired(authService))
		{
			matches.GET("", matchHandler.GetMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.GET("/:id/messages", messageHandler.GetMessages)
			matches.POST("/:id/messages", messageHandler.SendMessage)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthRequired(authService))
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		}
	}

	router.GET("/ws", middleware.AuthRequired(authService), wsHandler.HandleWebSocket)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
