package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yuliapopova/book_all/transaction/internal/domain"
)

type GormRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Deposit(params domain.DepositParams) (domain.TransactionDetails, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var details domain.TransactionDetails
	var entry domain.TransactionEntry

	amount := int64(params.Amount * 100)

	err = tx.QueryRow(
		`INSERT INTO transactions (user_id, amount, status, type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		params.UserID, amount, domain.StatusPending, domain.TypeDeposit,
	).Scan(&details.Transaction.ID, &details.Transaction.CreatedAt)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	details.Transaction.UserID = params.UserID
	details.Transaction.Amount = amount
	details.Transaction.Status = domain.StatusPending
	details.Transaction.Type = domain.TypeDeposit
	details.Transaction.UpdatedAt = details.Transaction.CreatedAt

	err = tx.QueryRow(
		`INSERT INTO transaction_entries (transaction_id, amount, type) VALUES ($1, $2, $3) RETURNING id`,
		details.Transaction.ID, amount, string(domain.TypeDeposit),
	).Scan(&entry.ID)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create transaction entry: %w", err)
	}

	entry.TransactionID = details.Transaction.ID
	entry.Amount = amount
	entry.Type = string(domain.TypeDeposit)

	details.Entries = append(details.Entries, entry)

	if err := tx.Commit(); err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return details, nil
}

func (r *GormRepository) Withdraw(params domain.WithdrawParams) (domain.TransactionDetails, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var details domain.TransactionDetails
	var entry domain.TransactionEntry

	amount := int64(params.Amount * 100)

	err = tx.QueryRow(
		`INSERT INTO transactions (account_id, amount, status, type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		params.AccountID, amount, domain.StatusPending, domain.TypeWithdraw,
	).Scan(&details.Transaction.ID, &details.Transaction.CreatedAt)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	details.Transaction.AccountID = params.AccountID
	details.Transaction.Amount = amount
	details.Transaction.Status = domain.StatusPending
	details.Transaction.Type = domain.TypeWithdraw
	details.Transaction.UpdatedAt = details.Transaction.CreatedAt

	err = tx.QueryRow(
		`INSERT INTO transaction_entries (transaction_id, amount, type) VALUES ($1, $2, $3) RETURNING id`,
		details.Transaction.ID, amount, string(domain.TypeWithdraw),
	).Scan(&entry.ID)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create transaction entry: %w", err)
	}

	entry.TransactionID = details.Transaction.ID
	entry.Amount = amount
	entry.Type = string(domain.TypeWithdraw)

	details.Entries = append(details.Entries, entry)

	if err := tx.Commit(); err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return details, nil
}

func (r *GormRepository) Transfer(params domain.TransferParams) (domain.TransactionDetails, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var details domain.TransactionDetails
	var entry1, entry2 domain.TransactionEntry

	amount := int64(params.Amount * 100)

	err = tx.QueryRow(
		`INSERT INTO transactions (user_id, amount, status, type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		params.UserID, amount, domain.StatusPending, domain.TypeTransfer,
	).Scan(&details.Transaction.ID, &details.Transaction.CreatedAt)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	details.Transaction.UserID = params.UserID
	details.Transaction.Amount = amount
	details.Transaction.Status = domain.StatusPending
	details.Transaction.Type = domain.TypeTransfer
	details.Transaction.UpdatedAt = details.Transaction.CreatedAt

	err = tx.QueryRow(
		`INSERT INTO transaction_entries (transaction_id, account_id, amount, type) VALUES ($1, $2, $3, $4) RETURNING id`,
		details.Transaction.ID, params.UserID, -amount, string(domain.TypeWithdraw),
	).Scan(&entry1.ID)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create debit entry: %w", err)
	}

	entry1.TransactionID = details.Transaction.ID
	entry1.Amount = -amount
	entry1.Type = string(domain.TypeWithdraw)

	err = tx.QueryRow(
		`INSERT INTO transaction_entries (transaction_id, account_id, amount, type) VALUES ($1, $2, $3, $4) RETURNING id`,
		details.Transaction.ID, params.Recipient, amount, string(domain.TypeDeposit),
	).Scan(&entry2.ID)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to create credit entry: %w", err)
	}

	entry2.TransactionID = details.Transaction.ID
	entry2.Amount = amount
	entry2.Type = string(domain.TypeDeposit)

	details.Entries = append(details.Entries, entry1, entry2)

	if err := tx.Commit(); err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return details, nil
}

func (r *GormRepository) UpdateTransactionStatus(id uint64, status domain.TransactionStatus) error {
	_, err := r.db.Exec(
		`UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}

func (r *GormRepository) GetTransactions(params domain.GetTransactionsParams) ([]domain.Transaction, error) {
	query := `SELECT id, user_id, amount, status, type, created_at, updated_at FROM transactions WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if params.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *params.UserID)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		var status, txType string
		if err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.Amount, &status, &txType, &tx.CreatedAt, &tx.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		tx.Status = domain.TransactionStatus(status)
		tx.Type = domain.TransactionType(txType)
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *GormRepository) GetTransactionDetails(id uint64) (domain.TransactionDetails, error) {
	var details domain.TransactionDetails

	err := r.db.QueryRow(
		`SELECT id, user_id, amount, status, type, created_at, updated_at FROM transactions WHERE id = $1`,
		id,
	).Scan(
		&details.Transaction.ID, &details.Transaction.UserID, &details.Transaction.Amount,
		&details.Transaction.Status, &details.Transaction.Type,
		&details.Transaction.CreatedAt, &details.Transaction.UpdatedAt,
	)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	rows, err := r.db.Query(
		`SELECT id, transaction_id, amount, type, created_at FROM transaction_entries WHERE transaction_id = $1`,
		id,
	)
	if err != nil {
		return domain.TransactionDetails{}, fmt.Errorf("failed to get transaction entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry domain.TransactionEntry
		if err := rows.Scan(
			&entry.ID, &entry.TransactionID, &entry.Amount, &entry.Type, &entry.CreatedAt,
		); err != nil {
			return domain.TransactionDetails{}, fmt.Errorf("failed to scan entry: %w", err)
		}
		details.Entries = append(details.Entries, entry)
	}

	return details, nil
}
