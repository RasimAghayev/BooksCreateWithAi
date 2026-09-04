package domain

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type AccountResponse struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	OperationID uint64 `json:"operation_id"`
	Result      bool   `json:"result"`
}

type AccountServiceClient interface {
	Deposit(userID uint64, amount float64) (bool, error)
	Withdraw(accountID uint64, amount float64) (bool, error)
	Transfer(userID uint64, recipientID uint64, amount float64) (bool, error)
}

type KafkaPublisher interface {
	Publish(topic string, key string, data []byte) error
}

type KafkaConsumer interface {
	Consume(topic string, handler func(key string, data []byte) error)
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
	Deposit(params DepositParams) (TransactionDetails, error)
	Withdraw(params WithdrawParams) (TransactionDetails, error)
	Transfer(params TransferParams) (TransactionDetails, error)
	UpdateTransactionStatus(id uint64, status TransactionStatus) error
	GetTransactions(params GetTransactionsParams) ([]Transaction, error)
	GetTransactionDetails(id uint64) (TransactionDetails, error)
}
