package repository

import (
	"database/sql"
	"fmt"

	"github.com/yuliapopova/book_all/account/internal/domain"
)

type GormRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(userID uint64, initialBalance int64) (domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(
		`INSERT INTO accounts (user_id, balance, currency, status) VALUES ($1, $2, $3, $4) RETURNING id, user_id, balance, currency, status, created_at, updated_at`,
		userID, initialBalance, "USD", "active",
	).Scan(
		&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.Status,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to create account: %w", err)
	}
	return account, nil
}

func (r *GormRepository) GetByID(id uint64) (domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(
		`SELECT id, user_id, balance, currency, status, created_at, updated_at FROM accounts WHERE id = $1`,
		id,
	).Scan(
		&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.Status,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (r *GormRepository) GetByUserID(userID uint64) (domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(
		`SELECT id, user_id, balance, currency, status, created_at, updated_at FROM accounts WHERE user_id = $1`,
		userID,
	).Scan(
		&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.Status,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to get account by user_id: %w", err)
	}
	return account, nil
}

func (r *GormRepository) UpdateBalance(id uint64, amount int64) (domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 RETURNING id, user_id, balance, currency, status, created_at, updated_at`,
		amount, id,
	).Scan(
		&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.Status,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to update balance: %w", err)
	}
	return account, nil
}

func (r *GormRepository) TransferBalance(fromID, toID uint64, amount int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var fromBalance int64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, fromID).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return fmt.Errorf("failed to debit sender: %w", err)
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return fmt.Errorf("failed to credit recipient: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transfer: %w", err)
	}

	return nil
}
