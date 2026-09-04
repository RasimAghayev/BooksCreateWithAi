package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	_ "github.com/lib/pq"
	accountpb "github.com/yuliapopova/book_all/contracts/account"
	"github.com/yuliapopova/book_all/account/internal/config"
	"github.com/yuliapopova/book_all/account/internal/repository"
	"github.com/yuliapopova/book_all/account/internal/service"
	"github.com/yuliapopova/book_all/account/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := config.InitLogger()
	cfg := config.New()

	logger.Info().Msg("account service starting")

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
		Topic: "transaction_response",
	}
	defer kafkaWriter.Close()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBrokerHost},
		Topic:   cfg.KafkaTransactionTopic,
		GroupID: cfg.KafkaConsumerGroup,
	})
	defer kafkaReader.Close()

	repo := repository.New(db)
	kafkaPublisher := domain.NewKafkaPublisherAdapter(kafkaWriter)
	svc := service.New(repo, kafkaPublisher, logger)
	grpcServer := transport.New(svc, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC_PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}

	grpcServerInstance := grpc.NewServer()
	accountpb.RegisterAccountServiceServer(grpcServerInstance, grpcServer)
	reflection.Register(grpcServerInstance)

	logger.Info().Str("port", cfg.GRPC_PORT).Msg("gRPC server starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			m, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Error().Err(err).Msg("failed to read kafka message")
				continue
			}

			logger.Info().Str("topic", m.Topic).Msg("received transaction request")

			if err := svc.HandleTransactionRequest(ctx, string(m.Key), m.Value); err != nil {
				logger.Error().Err(err).Msg("failed to handle transaction request")
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info().Msg("shutting down...")
		cancel()
		grpcServerInstance.GracefulStop()
	}()

	if err := grpcServerInstance.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve")
	}
}
