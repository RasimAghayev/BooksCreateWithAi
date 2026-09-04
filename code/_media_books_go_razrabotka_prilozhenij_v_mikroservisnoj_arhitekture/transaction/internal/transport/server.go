package transport

import (
	"context"

	"github.com/rs/zerolog"
	transactionpb "github.com/yuliapopova/book_all/contracts/transaction"
	"github.com/yuliapopova/book_all/transaction/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	transactionpb.UnimplementedTransactionServiceServer
	repo      domain.Repository
	accounts  domain.AccountServiceClient
	kafka     domain.KafkaPublisher
	logger    *zerolog.Logger
}

func New(repo domain.Repository, accounts domain.AccountServiceClient, kafka domain.KafkaPublisher, logger *zerolog.Logger) *Server {
	return &Server{
		repo:      repo,
		accounts:  accounts,
		kafka:     kafka,
		logger:    logger,
	}
}

func (s *Server) Deposit(ctx context.Context, req *transactionpb.DepositRequest) (*transactionpb.DepositResponse, error) {
	service := NewTransactionServiceWrapper(s.repo, s.accounts, s.kafka, s.logger)
	details, err := service.Deposit(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deposit: %v", err)
	}

	return &transactionpb.DepositResponse{
		Success:  true,
		Message:  "Transaction created successfully",
		Details:  toProtoDetails(details),
	}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *transactionpb.WithdrawRequest) (*transactionpb.WithdrawResponse, error) {
	service := NewTransactionServiceWrapper(s.repo, s.accounts, s.kafka, s.logger)
	details, err := service.Withdraw(ctx, req.AccountId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to withdraw: %v", err)
	}

	return &transactionpb.WithdrawResponse{
		Success:  true,
		Message:  "Transaction created successfully",
		Details:  toProtoDetails(details),
	}, nil
}

func (s *Server) Transfer(ctx context.Context, req *transactionpb.TransferRequest) (*transactionpb.TransferResponse, error) {
	service := NewTransactionServiceWrapper(s.repo, s.accounts, s.kafka, s.logger)
	details, err := service.Transfer(ctx, req.UserId, req.Recipient, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer: %v", err)
	}

	return &transactionpb.TransferResponse{
		Success:  true,
		Message:  "Transaction created successfully",
		Details:  toProtoDetails(details),
	}, nil
}

func (s *Server) GetTransactions(ctx context.Context, req *transactionpb.GetTransactionsRequest) (*transactionpb.GetTransactionsResponse, error) {
	params := domain.GetTransactionsParams{
		UserID: &req.UserId,
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}

	transactions, err := s.repo.GetTransactions(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	var protoTransactions []*transactionpb.TransactionDetails
	for _, tx := range transactions {
		details, err := s.repo.GetTransactionDetails(tx.ID)
		if err != nil {
			s.logger.Error().Err(err).Uint64("transaction_id", tx.ID).Msg("failed to get transaction details")
			continue
		}
		protoTransactions = append(protoTransactions, toProtoDetails(details))
	}

	return &transactionpb.GetTransactionsResponse{
		Success:      true,
		Transactions: protoTransactions,
	}, nil
}

func toProtoTransactionStatus(status domain.TransactionStatus) transactionpb.TransactionStatus {
	switch status {
	case domain.StatusCompleted:
		return transactionpb.TransactionStatus_COMPLETED
	case domain.StatusFailed:
		return transactionpb.TransactionStatus_FAILED
	default:
		return transactionpb.TransactionStatus_PENDING
	}
}

func toProtoTransactionType(txType domain.TransactionType) transactionpb.TransactionType {
	switch txType {
	case domain.TypeWithdraw:
		return transactionpb.TransactionType_WITHDRAW
	case domain.TypeTransfer:
		return transactionpb.TransactionType_TRANSFER
	default:
		return transactionpb.TransactionType_DEPOSIT
	}
}

func toProtoDetails(details domain.TransactionDetails) *transactionpb.TransactionDetails {
	tx := &transactionpb.Transaction{
		Id:        details.Transaction.ID,
		UserId:    details.Transaction.UserID,
		AccountId: details.Transaction.AccountID,
		Amount:    details.Transaction.Amount,
		Status:    toProtoTransactionStatus(details.Transaction.Status),
		Type:      toProtoTransactionType(details.Transaction.Type),
	}

	var entries []*transactionpb.TransactionEntry
	for _, e := range details.Entries {
		entries = append(entries, &transactionpb.TransactionEntry{
			Id:            e.ID,
			TransactionId: e.TransactionID,
			Amount:        e.Amount,
			Type:          e.Type,
		})
	}

	return &transactionpb.TransactionDetails{
		Transaction: tx,
		Entries:     entries,
	}
}

type TransactionServiceWrapper struct {
	repo      domain.Repository
	accounts  domain.AccountServiceClient
	kafka     domain.KafkaPublisher
	logger    *zerolog.Logger
}

func NewTransactionServiceWrapper(repo domain.Repository, accounts domain.AccountServiceClient, kafka domain.KafkaPublisher, logger *zerolog.Logger) *TransactionServiceWrapper {
	return &TransactionServiceWrapper{
		repo:      repo,
		accounts:  accounts,
		kafka:     kafka,
		logger:    logger,
	}
}

func (s *TransactionServiceWrapper) Deposit(ctx context.Context, userID uint64, amount float64) (domain.TransactionDetails, error) {
	details, err := s.repo.Deposit(domain.DepositParams{UserID: userID, Amount: amount})
	if err != nil {
		return domain.TransactionDetails{}, err
	}

	s.logger.Info().Uint64("transaction_id", details.Transaction.ID).Msg("transaction created")
	return details, nil
}

func (s *TransactionServiceWrapper) Withdraw(ctx context.Context, accountID uint64, amount float64) (domain.TransactionDetails, error) {
	details, err := s.repo.Withdraw(domain.WithdrawParams{AccountID: accountID, Amount: amount})
	if err != nil {
		return domain.TransactionDetails{}, err
	}

	s.logger.Info().Uint64("transaction_id", details.Transaction.ID).Msg("transaction created")
	return details, nil
}

func (s *TransactionServiceWrapper) Transfer(ctx context.Context, userID uint64, recipientID uint64, amount float64) (domain.TransactionDetails, error) {
	details, err := s.repo.Transfer(domain.TransferParams{UserID: userID, Recipient: recipientID, Amount: amount})
	if err != nil {
		return domain.TransactionDetails{}, err
	}

	s.logger.Info().Uint64("transaction_id", details.Transaction.ID).Msg("transaction created")
	return details, nil
}
