package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DBDSN         string
	JWTSecret     string
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
	GRPC_PORT     string
}

func New() *Config {
	return &Config{
		DBDSN:            os.Getenv("DB_DSN"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		AccessTokenTTL:   getEnvAsDuration("ACCESS_TOKEN_TTL_MINUTES", 60*time.Minute),
		RefreshTokenTTL:  getEnvAsDuration("REFRESH_TOKEN_TTL_DAYS", 720*time.Hour),
		GRPC_PORT:        getEnv("GRPC_PORT", "50052"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value + "m"); err == nil {
			return dur
		}
	}
	return defaultValue
}

func InitLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger
}
