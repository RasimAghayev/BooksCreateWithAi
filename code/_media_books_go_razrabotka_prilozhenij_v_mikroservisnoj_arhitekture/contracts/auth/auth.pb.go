package auth

import "context"

type RegisterRequest struct {
	Email    string
	Password string
	Name     string
}

type RegisterResponse struct {
	Success bool
	Message string
	User    *User
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Success      bool
	Message      string
	AccessToken  string
	RefreshToken string
}

type LogoutRequest struct {
	RefreshToken string
}

type LogoutResponse struct {
	Success bool
}

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	Success      bool
	Message      string
	AccessToken  string
	RefreshToken string
}

type GetCurrentUserRequest struct {
	AccessToken string
}

type GetCurrentUserResponse struct {
	Success bool
	User    *User
}

type User struct {
	Id    uint64
	Email string
	Name  string
}

type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	GetCurrentUser(context.Context, *GetCurrentUserRequest) (*GetCurrentUserResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

type UnimplementedAuthServiceServer struct{}

func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}
func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, nil
}
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
func (UnimplementedAuthServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, nil
}
func (UnimplementedAuthServiceServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, nil
}
func (UnimplementedAuthServiceServer) GetCurrentUser(context.Context, *GetCurrentUserRequest) (*GetCurrentUserResponse, error) {
	return nil, nil
}

func RegisterAuthServiceServer(s interface{}, srv AuthServiceServer) {}
func NewAuthServiceClient(cc interface{}) AuthServiceClient {
	return &authServiceClient{cc}
}

type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts interface{}) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts interface{}) (*LoginResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts interface{}) (*RefreshResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts interface{}) (*LogoutResponse, error)
	GetCurrentUser(ctx context.Context, in *GetCurrentUserRequest, opts interface{}) (*GetCurrentUserResponse, error)
}

type authServiceClient struct {
	cc interface{}
}

func (c *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts interface{}) (*RegisterResponse, error) {
	return nil, nil
}
func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts interface{}) (*LoginResponse, error) {
	return nil, nil
}
func (c *authServiceClient) Refresh(ctx context.Context, in *RefreshRequest, opts interface{}) (*RefreshResponse, error) {
	return nil, nil
}
func (c *authServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts interface{}) (*LogoutResponse, error) {
	return nil, nil
}
func (c *authServiceClient) GetCurrentUser(ctx context.Context, in *GetCurrentUserRequest, opts interface{}) (*GetCurrentUserResponse, error) {
	return nil, nil
}