package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/yuliapopova/book_all/auth/internal/domain"
)

var (
	JWT_SECRET              = getEnv("JWT_SECRET", "secret")
	ACCESS_TOKEN_TTL        = 60 * time.Minute
	REFRESH_TOKEN_TTL       = 720 * time.Hour
)

type AuthService struct {
	repo   domain.Repository
	logger *zerolog.Logger
}

func New(repo domain.Repository, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (domain.User, error) {
	if email == "" || password == "" {
		return domain.User{}, errors.New("email and password are required")
	}

	_, err := s.repo.GetByEmail(email)
	if err == nil {
		return domain.User{}, errors.New("user with this email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(email, string(passwordHash), name)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info().Uint64("user_id", user.ID).Msg("user registered")
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (domain.Session, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return domain.Session{}, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.Session{}, errors.New("invalid email or password")
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.repo.CreateRefreshToken(user.ID, refreshToken, time.Now().Add(REFRESH_TOKEN_TTL)); err != nil {
		return domain.Session{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	s.logger.Info().Uint64("user_id", user.ID).Msg("user logged in")

	return domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	rt, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		return "", errors.New("refresh token expired")
	}

	user, err := s.repo.GetByID(rt.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.repo.DeleteRefreshToken(refreshToken); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	s.logger.Info().Msg("user logged out")
	return nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, accessToken string) (domain.User, error) {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return domain.User{}, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := s.repo.GetByID(claims.UserID)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *AuthService) generateAccessToken(userID uint64, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iss":     "auth-service",
		"exp":     time.Now().Add(ACCESS_TOKEN_TTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func (s *AuthService) generateRefreshToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"is_refresh": true,
		"iss":        "auth-service",
		"exp":        time.Now().Add(REFRESH_TOKEN_TTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func (s *AuthService) parseAccessToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &domain.JWTClaims{
			UserID:   uint64(claims["user_id"].(float64)),
			Email:    claims["email"].(string),
			IsRefresh: claims["is_refresh"].(bool),
			Iss:      claims["iss"].(string),
			Exp:      int64(claims["exp"].(float64)),
			Iat:      int64(claims["iat"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token claims")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var _ = rsa.ParsePKCS1PrivateKey
