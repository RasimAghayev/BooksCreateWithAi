package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/yuliapopova/book_all/transaction/internal/domain"
)

type TransactionService struct {
	repo      domain.Repository
	accounts  domain.AccountServiceClient
	kafka     domain.KafkaPublisher
	logger    *zerolog.Logger
}

func New(repo domain.Repository, accounts domain.AccountServiceClient, kafka domain.KafkaPublisher, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{
		repo:      repo,
		accounts:  accounts,
		kafka:     kafka,
		logger:    logger,
	}
}

func (s *TransactionService) Deposit(ctx context.Context, userID uint64, amount float64) (domain.TransactionDetails, error) {
	params := domain.DepositParams{
		UserID: userID,
		Amount: amount,
	}

	details, err := s.repo.Deposit(params)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to deposit: %w", err)
	}

	data, err := json.Marshal(details.Transaction)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal transaction")
	} else {
		if err := s.kafka.Publish("transaction_data", fmt.Sprintf("%d", details.Transaction.ID), data); err != nil {
			s.logger.Error().Err(err).Msg("failed to publish transaction event")
			if err := s.repo.UpdateTransactionStatus(details.Transaction.ID, domain.StatusFailed); err != nil {
				s.logger.Error().Err(err).Msg("failed to update transaction status")
			}
			return domain.TransactionDetails{}, fmt.Errorf("failed to publish transaction event: %w", err)
		}
	}

	return details, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID uint64, amount float64) (domain.TransactionDetails, error) {
	params := domain.WithdrawParams{
		AccountID: accountID,
		Amount:    amount,
	}

	details, err := s.repo.Withdraw(params)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to withdraw: %w", err)
	}

	data, err := json.Marshal(details.Transaction)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal transaction")
	} else {
		if err := s.kafka.Publish("transaction_data", fmt.Sprintf("%d", details.Transaction.ID), data); err != nil {
			s.logger.Error().Err(err).Msg("failed to publish transaction event")
			if err := s.repo.UpdateTransactionStatus(details.Transaction.ID, domain.StatusFailed); err != nil {
				s.logger.Error().Err(err).Msg("failed to update transaction status")
			}
			return domain.TransactionDetails{}, fmt.Errorf("failed to publish transaction event: %w", err)
		}
	}

	return details, nil
}

func (s *TransactionService) Transfer(ctx context.Context, userID uint64, recipientID uint64, amount float64) (domain.TransactionDetails, error) {
	params := domain.TransferParams{
		UserID:    userID,
		Recipient: recipientID,
		Amount:    amount,
	}

	details, err := s.repo.Transfer(params)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to transfer: %w", err)
	}

	data, err := json.Marshal(details.Transaction)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal transaction")
	} else {
		if err := s.kafka.Publish("transaction_data", fmt.Sprintf("%d", details.Transaction.ID), data); err != nil {
			s.logger.Error().Err(err).Msg("failed to publish transaction event")
			if err := s.repo.UpdateTransactionStatus(details.Transaction.ID, domain.StatusFailed); err != nil {
				s.logger.Error().Err(err).Msg("failed to update transaction status")
			}
			return domain.TransactionDetails{}, fmt.Errorf("failed to publish transaction event: %w", err)
		}
	}

	return details, nil
}

func (s *TransactionService) HandleAccountResponse(ctx context.Context, topic string, key string, data []byte) error {
	var response domain.AccountResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	if response.RequestType == "" {
		return fmt.Errorf("missing request type")
	}
	if response.UserID == 0 {
		return fmt.Errorf("missing user ID")
	}
	if response.OperationID == 0 {
		return fmt.Errorf("missing operation ID")
	}

	status := domain.StatusFailed
	if response.Result {
		status = domain.StatusCompleted
	}

	if err := s.repo.UpdateTransactionStatus(response.OperationID, status); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	s.logger.Info().
		Uint64("transaction_id", response.OperationID).
		Str("status", string(status)).
		Msg("transaction status updated")

	return nil
}

func (s *TransactionService) GetTransactionsWithDetails(ctx context.Context, params domain.GetTransactionsParams) ([]domain.TransactionDetails, error) {
	transactions, err := s.repo.GetTransactions(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var detailsList []domain.TransactionDetails
	for _, tx := range transactions {
		details, err := s.repo.GetTransactionDetails(tx.ID)
		if err != nil {
			s.logger.Error().Err(err).Uint64("transaction_id", tx.ID).Msg("failed to get transaction details")
			continue
		}
		detailsList = append(detailsList, details)
	}

	return detailsList, nil
}
