package domain

type TransactionRequest struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	AccountID   uint64 `json:"account_id"`
	RecipientID uint64 `json:"recipient_id"`
	Amount      int64  `json:"amount"`
	OperationID uint64 `json:"operation_id"`
	Timestamp   int64  `json:"timestamp"`
}
