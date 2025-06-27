package models

import (
    "fmt"
    "time"
)

// TransactionType represents allowed transaction types.
const (
    TransactionDeposit    = "deposit"
    TransactionWithdrawal = "withdrawal"
    TransactionBuy        = "buy"
    TransactionSell       = "sell"
    TransactionTransfer   = "transfer"
)

// TransactionStatus indicates the state of a transaction.
const (
    StatusPending = "pending"
    StatusSuccess = "success"
    StatusFailed  = "failed"
)

// Transaction represents a money movement.
type Transaction struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    WalletID  string    `bson:"wallet_id" json:"wallet_id"`
    Amount    float64   `bson:"amount" json:"amount"`
    Type      string    `bson:"type" json:"type"`
    Status    string    `bson:"status" json:"status"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

var (
    typeSet = map[string]struct{}{
        TransactionDeposit:    {},
        TransactionWithdrawal: {},
        TransactionBuy:        {},
        TransactionSell:       {},
        TransactionTransfer:   {},
    }
    statusSet = map[string]struct{}{
        StatusPending: {},
        StatusSuccess: {},
        StatusFailed:  {},
    }
)

// Validate verifies the transaction fields.
func (t *Transaction) Validate() error {
    if t.WalletID == "" {
        return fmt.Errorf("wallet_id cannot be empty")
    }
    if t.Amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if _, ok := typeSet[t.Type]; !ok {
        return fmt.Errorf("invalid transaction type: %s", t.Type)
    }
    if _, ok := statusSet[t.Status]; !ok {
        return fmt.Errorf("invalid transaction status: %s", t.Status)
    }
    return nil
}
