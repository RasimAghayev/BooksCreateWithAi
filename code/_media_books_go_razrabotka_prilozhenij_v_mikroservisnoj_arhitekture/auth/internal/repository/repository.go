package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yuliapopova/book_all/auth/internal/domain"
)

type GormRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateUser(email, passwordHash, name string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, email, password_hash, name, created_at, updated_at`,
		email, passwordHash, name,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r *GormRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *GormRepository) GetByID(id uint64) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *GormRepository) CreateRefreshToken(userID uint64, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *GormRepository) GetRefreshToken(token string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.QueryRow(
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`,
		token,
	).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt,
	)
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return rt, nil
}

func (r *GormRepository) DeleteRefreshToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}
