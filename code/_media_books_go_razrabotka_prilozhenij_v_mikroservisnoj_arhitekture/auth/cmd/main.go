package main

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	authpb "github.com/yuliapopova/book_all/contracts/auth"
	"github.com/yuliapopova/book_all/auth/internal/config"
	"github.com/yuliapopova/book_all/auth/internal/repository"
	"github.com/yuliapopova/book_all/auth/internal/service"
	"github.com/yuliapopova/book_all/auth/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := config.InitLogger()
	cfg := config.New()

	logger.Info().Msg("auth service starting")

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}
	logger.Info().Msg("connected to database")

	repo := repository.New(db)
	svc := service.New(repo, logger)
	server := transport.New(svc, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC_PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	logger.Info().Str("port", cfg.GRPC_PORT).Msg("gRPC server starting")

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve")
	}
}

var _ = zerolog.TimeFieldFormat
