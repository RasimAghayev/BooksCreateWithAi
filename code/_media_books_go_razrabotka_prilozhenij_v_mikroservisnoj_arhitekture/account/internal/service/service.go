package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/yuliapopova/book_all/account/internal/domain"
	accountpb "github.com/yuliapopova/book_all/contracts/account"
)

type AccountService struct {
	repo  domain.Repository
	kafka domain.KafkaPublisher
	logger *zerolog.Logger
}

func New(repo domain.Repository, kafka domain.KafkaPublisher, logger *zerolog.Logger) *AccountService {
	return &AccountService{
		repo:   repo,
		kafka:  kafka,
		logger: logger,
	}
}

func (s *AccountService) Deposit(ctx context.Context, userID uint64, amount float64) (bool, error) {
	cents := int64(amount * 100)

	account, err := s.repo.GetByUserID(userID)
	if err != nil {
		account, err = s.repo.Create(userID, cents)
		if err != nil {
			return false, fmt.Errorf("failed to create account: %w", err)
		}
	}

	account, err = s.repo.UpdateBalance(account.ID, cents)
	if err != nil {
		return false, fmt.Errorf("failed to deposit: %w", err)
	}

	s.logger.Info().
		Uint64("account_id", account.ID).
		Int64("new_balance", account.Balance).
		Msg("deposit completed")

	return true, nil
}

func (s *AccountService) Withdraw(ctx context.Context, accountID uint64, amount float64) (bool, error) {
	cents := int64(amount * 100)

	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return false, fmt.Errorf("failed to get account: %w", err)
	}

	if account.Balance < cents {
		return false, fmt.Errorf("insufficient funds")
	}

	account, err = s.repo.UpdateBalance(accountID, -cents)
	if err != nil {
		return false, fmt.Errorf("failed to withdraw: %w", err)
	}

	s.logger.Info().
		Uint64("account_id", account.ID).
		Int64("new_balance", account.Balance).
		Msg("withdraw completed")

	return true, nil
}

func (s *AccountService) Transfer(ctx context.Context, fromUserID uint64, toUserID uint64, amount float64) (bool, error) {
	cents := int64(amount * 100)

	fromAccount, err := s.repo.GetByUserID(fromUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get sender account: %w", err)
	}

	toAccount, err := s.repo.GetByUserID(toUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get recipient account: %w", err)
	}

	if fromAccount.Balance < cents {
		return false, fmt.Errorf("insufficient funds")
	}

	err = s.repo.TransferBalance(fromAccount.ID, toAccount.ID, cents)
	if err != nil {
		return false, fmt.Errorf("failed to transfer: %w", err)
	}

	s.logger.Info().
		Uint64("from_account", fromAccount.ID).
		Uint64("to_account", toAccount.ID).
		Int64("amount", cents).
		Msg("transfer completed")

	return true, nil
}

func (s *AccountService) GetBalance(ctx context.Context, accountID uint64) (float64, error) {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}

	return float64(account.Balance) / 100.0, nil
}

func (s *AccountService) HandleTransactionRequest(ctx context.Context, key string, data []byte) error {
	var txRequest domain.TransactionRequest
	if err := json.Unmarshal(data, &txRequest); err != nil {
		return fmt.Errorf("failed to unmarshal transaction request: %w", err)
	}

	var result bool
	var err error

	switch txRequest.RequestType {
	case "deposit":
		result, err = s.Deposit(ctx, txRequest.UserID, float64(txRequest.Amount)/100.0)
	case "withdraw":
		result, err = s.Withdraw(ctx, txRequest.AccountID, float64(txRequest.Amount)/100.0)
	case "transfer":
		result, err = s.Transfer(ctx, txRequest.UserID, txRequest.RecipientID, float64(txRequest.Amount)/100.0)
	default:
		return fmt.Errorf("unknown request type: %s", txRequest.RequestType)
	}

	if err != nil {
		result = false
	}

	response := domain.AccountResponse{
		RequestType: txRequest.RequestType,
		UserID:      txRequest.UserID,
		OperationID: txRequest.OperationID,
		Result:      result,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if err := s.kafka.Publish("transaction_response", fmt.Sprintf("%d", txRequest.OperationID), responseData); err != nil {
		return fmt.Errorf("failed to publish response: %w", err)
	}

	s.logger.Info().
		Uint64("operation_id", txRequest.OperationID).
		Str("request_type", txRequest.RequestType).
		Bool("result", result).
		Msg("transaction request processed")

	return nil
}

func (s *AccountService) GetBalanceGRPC(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.GetBalanceResponse, error) {
	account, err := s.repo.GetByID(req.AccountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &accountpb.GetBalanceResponse{
		Success: true,
		Balance: float64(account.Balance) / 100.0,
	}, nil
}

type TransactionRequest struct {
	RequestType  string `json:"request_type"`
	UserID       uint64 `json:"user_id"`
	AccountID    uint64 `json:"account_id"`
	RecipientID  uint64 `json:"recipient_id"`
	Amount       int64  `json:"amount"`
	OperationID  uint64 `json:"operation_id"`
	Timestamp    int64  `json:"timestamp"`
}

func init() {
	_ = time.Now
}
