package domain

import "time"

type User struct {
	ID           uint64
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type JWTClaims struct {
	UserID   uint64 `json:"user_id"`
	Email    string `json:"email"`
	IsRefresh bool   `json:"is_refresh"`
	Iss       string `json:"iss"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}

type Session struct {
	AccessToken  string
	RefreshToken string
}

type Repository interface {
	CreateUser(email, passwordHash, name string) (User, error)
	GetByEmail(email string) (User, error)
	GetByID(id uint64) (User, error)
	CreateRefreshToken(userID uint64, token string, expiresAt time.Time) error
	GetRefreshToken(token string) (RefreshToken, error)
	DeleteRefreshToken(token string) error
}
