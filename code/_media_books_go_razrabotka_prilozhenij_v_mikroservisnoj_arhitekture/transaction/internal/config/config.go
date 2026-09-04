package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBDSN              string
	AccountGRPCHost    string
	KafkaBrokerHost    string
	KafkaConsumerGroup string
	KafkaTransactionTopic string
	GRPC_PORT          string
}

func New() *Config {
	return &Config{
		DBDSN:               os.Getenv("DB_DSN"),
		AccountGRPCHost:     getEnv("ACCOUNT_GRPC_HOST", "localhost:50051"),
		KafkaBrokerHost:     getEnv("KAFKA_BROKER_HOST", "localhost:9092"),
		KafkaConsumerGroup:  getEnv("KAFKA_CONSUMER_GROUP", "consumer"),
		KafkaTransactionTopic: getEnv("KAFKA_TRANSACTION_TOPIC", "transaction_response"),
		GRPC_PORT:           getEnv("GRPC_PORT", "50054"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
	}
	return defaultValue
}
