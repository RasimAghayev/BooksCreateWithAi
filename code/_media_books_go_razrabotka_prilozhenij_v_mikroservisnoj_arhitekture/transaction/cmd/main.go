package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	_ "github.com/lib/pq"
	transactionpb "github.com/yuliapopova/book_all/contracts/transaction"
	"github.com/yuliapopova/book_all/transaction/internal/config"
	"github.com/yuliapopova/book_all/transaction/internal/domain"
	"github.com/yuliapopova/book_all/transaction/internal/repository"
	"github.com/yuliapopova/book_all/transaction/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFieldFormat
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.New()
	logger.Info().Msg("configuration loaded")

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}
	logger.Info().Msg("connected to database")

	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP(cfg.KafkaBrokerHost),
		Topic: cfg.KafkaTransactionTopic,
	}
	defer kafkaWriter.Close()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBrokerHost},
		Topic:   "transaction_response",
		GroupID: cfg.KafkaConsumerGroup,
	})
	defer kafkaReader.Close()

	repo := repository.New(db)
	kafkaPublisher := domain.NewKafkaPublisherAdapter(kafkaWriter)
	server := transport.New(repo, nil, kafkaPublisher, &logger)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC_PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}

	grpcServer := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	logger.Info().Str("port", cfg.GRPC_PORT).Msg("gRPC server starting")

	go consumeKafka(kafkaReader, server, &logger)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve")
	}
}

func consumeKafka(reader *kafka.Reader, server *transport.Server, logger *zerolog.Logger) {
	for {
		m, err := reader.ReadMessage(nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to read message")
			continue
		}

		logger.Info().
			Str("topic", m.Topic).
			Str("partition", fmt.Sprintf("%d", m.Partition)).
			Str("offset", fmt.Sprintf("%d", m.Offset)).
			Msg("received message")

		m.Value = m.Value
		_ = m
	}
}
