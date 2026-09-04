package transport

import (
	"context"

	"github.com/rs/zerolog"
	accountpb "github.com/yuliapopova/book_all/contracts/account"
	"github.com/yuliapopova/book_all/account/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	accountpb.UnimplementedAccountServiceServer
	service *service.AccountService
	logger  *zerolog.Logger
}

func New(svc *service.AccountService, logger *zerolog.Logger) *Server {
	return &Server{
		service: svc,
		logger:  logger,
	}
}

func (s *Server) Deposit(ctx context.Context, req *accountpb.DepositRequest) (*accountpb.DepositResponse, error) {
	result, err := s.service.Deposit(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deposit: %v", err)
	}

	return &accountpb.DepositResponse{
		Success:  result,
		Message:  "Deposit completed",
	}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *accountpb.WithdrawRequest) (*accountpb.WithdrawResponse, error) {
	result, err := s.service.Withdraw(ctx, req.AccountId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to withdraw: %v", err)
	}

	return &accountpb.WithdrawResponse{
		Success:  result,
		Message:  "Withdraw completed",
	}, nil
}

func (s *Server) Transfer(ctx context.Context, req *accountpb.TransferRequest) (*accountpb.TransferResponse, error) {
	result, err := s.service.Transfer(ctx, req.AccountId, req.RecipientId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer: %v", err)
	}

	return &accountpb.TransferResponse{
		Success:  result,
		Message:  "Transfer completed",
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.GetBalanceResponse, error) {
	balance, err := s.service.GetBalance(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	return &accountpb.GetBalanceResponse{
		Success: true,
		Balance: balance,
	}, nil
}
