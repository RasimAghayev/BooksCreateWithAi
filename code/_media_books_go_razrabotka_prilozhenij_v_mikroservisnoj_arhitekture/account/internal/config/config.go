package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DBDSN              string
	KafkaBrokerHost    string
	KafkaConsumerGroup string
	KafkaTransactionTopic string
	GRPC_PORT          string
}

func New() *Config {
	return &Config{
		DBDSN:                os.Getenv("DB_DSN"),
		KafkaBrokerHost:      getEnv("KAFKA_BROKER_HOST", "localhost:9092"),
		KafkaConsumerGroup:   getEnv("KAFKA_CONSUMER_GROUP", "consumer-account"),
		KafkaTransactionTopic: getEnv("KAFKA_TRANSACTION_TOPIC", "transaction_request"),
		GRPC_PORT:            getEnv("GRPC_PORT", "50051"),
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
