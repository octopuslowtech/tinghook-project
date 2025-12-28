package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	CORSOrigins string

	DatabaseURL string
	RedisURL    string

	JWTSecret     string
	JWTExpiration string

	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),

		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
