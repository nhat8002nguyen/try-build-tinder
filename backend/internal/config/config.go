package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	JWTExpireHours int

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	FacebookClientID     string
	FacebookClientSecret string
	FacebookRedirectURL  string

	StorageType     string // "local" or "s3"
	LocalStorageDir string
	S3Bucket        string
	S3Region        string
	AWSAccessKey    string
	AWSSecretKey    string

	Environment string
}

func Load() *Config {
	godotenv.Load()

	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tinder_clone?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpireHours: jwtExpire,

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/oauth/google/callback"),

		FacebookClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
		FacebookRedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", "http://localhost:8080/api/auth/oauth/facebook/callback"),

		StorageType:     getEnv("STORAGE_TYPE", "local"),
		LocalStorageDir: getEnv("LOCAL_STORAGE_DIR", "./uploads"),
		S3Bucket:        getEnv("S3_BUCKET", ""),
		S3Region:        getEnv("S3_REGION", "us-east-1"),
		AWSAccessKey:    getEnv("AWS_ACCESS_KEY", ""),
		AWSSecretKey:    getEnv("AWS_SECRET_KEY", ""),

		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
