package transaction

import (
	context "context"

	empty "google.golang.org/protobuf/types/known/emptypb"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionStatus int32

const (
	TransactionStatus_PENDING   TransactionStatus = 0
	TransactionStatus_COMPLETED TransactionStatus = 1
	TransactionStatus_FAILED    TransactionStatus = 2
)

type TransactionType int32

const (
	TransactionType_DEPOSIT  TransactionType = 0
	TransactionType_WITHDRAW TransactionType = 1
	TransactionType_TRANSFER TransactionType = 2
)

type DepositRequest struct {
	UserId uint64
	Amount float64
}

type DepositResponse struct {
	Success bool
	Message string
	Details *TransactionDetails
}

type WithdrawRequest struct {
	AccountId uint64
	Amount    float64
}

type WithdrawResponse struct {
	Success bool
	Message string
	Details *TransactionDetails
}

type TransferRequest struct {
	UserId    uint64
	Recipient uint64
	Amount    float64
}

type TransferResponse struct {
	Success bool
	Message string
	Details *TransactionDetails
}

type GetTransactionsRequest struct {
	UserId uint64
	Limit  uint32
	Offset uint32
}

type GetTransactionsResponse struct {
	Success      bool
	Transactions []*TransactionDetails
}

type TransactionEntry struct {
	Id            uint64
	TransactionId uint64
	Amount        int64
	Type          string
}

type Transaction struct {
	Id        uint64
	UserId    uint64
	AccountId uint64
	Amount    int64
	Status    TransactionStatus
	Type      TransactionType
	CreatedAt *timestamp.Timestamp
	UpdatedAt *timestamp.Timestamp
}

type TransactionDetails struct {
	Transaction *Transaction
	Entries     []*TransactionEntry
}

type TransactionServiceServer interface {
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
	GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error)
	HandleAccountResponse(context.Context, *empty.Empty) (*empty.Empty, error)
	mustEmbedUnimplementedTransactionServiceServer()
}

type UnimplementedTransactionServiceServer struct{}

func (UnimplementedTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {}
func (UnimplementedTransactionServiceServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, nil
}
func (UnimplementedTransactionServiceServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, nil
}
func (UnimplementedTransactionServiceServer) Transfer(context.Context, *TransferRequest) (*TransferResponse, error) {
	return nil, nil
}
func (UnimplementedTransactionServiceServer) GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	return nil, nil
}
func (UnimplementedTransactionServiceServer) HandleAccountResponse(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, nil
}

func RegisterTransactionServiceServer(s interface{}, srv TransactionServiceServer) {
	_ = s
	_ = srv
}

var _ = empty.Empty{}
var _ = timestamp.Timestamp{}