package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AccountGRPCHost     string
	AuthGRPCHost        string
	TransactionGRPCHost string
	JWTSecret           string
	HTTP_PORT           string
}

func New() *Config {
	return &Config{
		AccountGRPCHost:     getEnv("ACCOUNT_GRPC_HOST", "localhost:50051"),
		AuthGRPCHost:        getEnv("AUTH_GRPC_HOST", "localhost:50052"),
		TransactionGRPCHost: getEnv("TRANSACTION_GRPC_HOST", "localhost:50054"),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		HTTP_PORT:           getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger
}
