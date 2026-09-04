package transport

import (
	"context"

	"github.com/rs/zerolog"
	authpb "github.com/yuliapopova/book_all/contracts/auth"
	"github.com/yuliapopova/book_all/auth/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	service *service.AuthService
	logger  *zerolog.Logger
}

func New(svc *service.AuthService, logger *zerolog.Logger) *Server {
	return &Server{
		service: svc,
		logger:  logger,
	}
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user, err := s.service.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &authpb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		User: &authpb.User{
			Id:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	session, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login: %v", err)
	}

	return &authpb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if err := s.service.Logout(ctx, req.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout: %v", err)
	}

	return &authpb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	accessToken, err := s.service.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to refresh: %v", err)
	}

	return &authpb.RefreshResponse{
		Success:      true,
		AccessToken:  accessToken,
	}, nil
}

func (s *Server) GetCurrentUser(ctx context.Context, req *authpb.GetCurrentUserRequest) (*authpb.GetCurrentUserResponse, error) {
	user, err := s.service.GetCurrentUser(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get user: %v", err)
	}

	return &authpb.GetCurrentUserResponse{
		Success: true,
		User: &authpb.User{
			Id:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
