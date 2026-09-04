package domain

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Account struct {
	ID        uint64
	UserID    uint64
	Balance   int64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountResponse struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	OperationID uint64 `json:"operation_id"`
	Result      bool   `json:"result"`
}

type KafkaPublisher interface {
	Publish(topic string, key string, data []byte) error
}

type KafkaPublisherAdapter struct {
	writer *kafka.Writer
}

func NewKafkaPublisherAdapter(writer *kafka.Writer) *KafkaPublisherAdapter {
	return &KafkaPublisherAdapter{writer: writer}
}

func (a *KafkaPublisherAdapter) Publish(topic string, key string, data []byte) error {
	return a.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})
}

type Repository interface {
	Create(userID uint64, initialBalance int64) (Account, error)
	GetByID(id uint64) (Account, error)
	GetByUserID(userID uint64) (Account, error)
	UpdateBalance(id uint64, amount int64) (Account, error)
	TransferBalance(fromID, toID uint64, amount int64) error
}
