package account

import "context"

type DepositRequest struct {
	UserId uint64
	Amount float64
}

type DepositResponse struct {
	Success   bool
	Message   string
	AccountId uint64
	Balance   float64
}

type WithdrawRequest struct {
	AccountId uint64
	Amount    float64
}

type WithdrawResponse struct {
	Success bool
	Message string
	Balance float64
}

type TransferRequest struct {
	AccountId   uint64
	RecipientId uint64
	Amount      float64
}

type TransferResponse struct {
	Success          bool
	Message          string
	SenderBalance    float64
	RecipientBalance float64
}

type GetBalanceRequest struct {
	AccountId uint64
}

type GetBalanceResponse struct {
	Success bool
	Balance float64
}

type AccountResponse struct {
	RequestType string
	UserId      uint64
	OperationID uint64
	Result      bool
}

type AccountServiceServer interface {
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

type UnimplementedAccountServiceServer struct{}

func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}
func (UnimplementedAccountServiceServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, nil
}
func (UnimplementedAccountServiceServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, nil
}
func (UnimplementedAccountServiceServer) Transfer(context.Context, *TransferRequest) (*TransferResponse, error) {
	return nil, nil
}
func (UnimplementedAccountServiceServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, nil
}

func RegisterAccountServiceServer(s interface{}, srv AccountServiceServer) {}
func NewAccountServiceClient(cc interface{}) AccountServiceClient {
	return &accountServiceClient{cc}
}

type AccountServiceClient interface {
	Deposit(ctx context.Context, in *DepositRequest, opts interface{}) (*DepositResponse, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts interface{}) (*WithdrawResponse, error)
	Transfer(ctx context.Context, in *TransferRequest, opts interface{}) (*TransferResponse, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts interface{}) (*GetBalanceResponse, error)
}

type accountServiceClient struct {
	cc interface{}
}

func (c *accountServiceClient) Deposit(ctx context.Context, in *DepositRequest, opts interface{}) (*DepositResponse, error) {
	return nil, nil
}
func (c *accountServiceClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts interface{}) (*WithdrawResponse, error) {
	return nil, nil
}
func (c *accountServiceClient) Transfer(ctx context.Context, in *TransferRequest, opts interface{}) (*TransferResponse, error) {
	return nil, nil
}
func (c *accountServiceClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts interface{}) (*GetBalanceResponse, error) {
	return nil, nil
}