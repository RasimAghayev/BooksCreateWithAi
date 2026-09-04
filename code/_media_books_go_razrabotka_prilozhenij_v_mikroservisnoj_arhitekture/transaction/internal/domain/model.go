package domain

import "time"

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
)

type TransactionType string

const (
	TypeDeposit  TransactionType = "deposit"
	TypeWithdraw TransactionType = "withdraw"
	TypeTransfer TransactionType = "transfer"
)

type TransactionEntry struct {
	ID           uint64
	TransactionID uint64
	Amount       int64
	Type         string
	CreatedAt    time.Time
}

type Transaction struct {
	ID        uint64
	UserID    uint64
	AccountID uint64
	Amount    int64
	Status    TransactionStatus
	Type      TransactionType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionDetails struct {
	Transaction Transaction
	Entries     []TransactionEntry
}

type DepositParams struct {
	UserID uint64
	Amount float64
}

type WithdrawParams struct {
	AccountID uint64
	Amount    float64
}

type TransferParams struct {
	UserID    uint64
	Recipient uint64
	Amount    float64
}

type GetTransactionsParams struct {
	UserID *uint64
	Limit  int32
	Offset int32
}
