package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/tinder-clone/backend/internal/config"
	"github.com/tinder-clone/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func New(cfg *config.Config) (*Database, error) {
	gormConfig := &gorm.Config{}
	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Warn().Err(err).Msg("Redis connection failed, continuing without cache")
		redisClient = nil
	}

	return &Database{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func (d *Database) AutoMigrate() error {
	// Check if migrations have already been run
	if d.DB.Migrator().HasTable(&models.User{}) {
		log.Info().Msg("Tables already exist, skipping migrations")
		return nil
	}

	return d.DB.AutoMigrate(
		&models.User{},
		&models.UserPhoto{},
		&models.OAuthAccount{},
		&models.Swipe{},
		&models.Match{},
		&models.Message{},
		&models.Notification{},
	)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	if d.Redis != nil {
		d.Redis.Close()
	}

	return sqlDB.Close()
}
