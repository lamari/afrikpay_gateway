package models

import (
    "fmt"
    "time"
)

// Wallet represents a crypto wallet owned by a user.
type Wallet struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    UserID    string    `bson:"user_id" json:"user_id"`
    Balance   float64   `bson:"balance" json:"balance"`
    Currency  string    `bson:"currency" json:"currency"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

var currencySet = map[string]struct{}{
    "USDT": {},
    "BTC":  {},
}

// Validate ensures the wallet fields are consistent.
func (w *Wallet) Validate() error {
    if w.UserID == "" {
        return fmt.Errorf("user_id cannot be empty")
    }
    if w.Balance < 0 {
        return fmt.Errorf("balance cannot be negative")
    }
    if _, ok := currencySet[w.Currency]; !ok {
        return fmt.Errorf("unsupported currency: %s", w.Currency)
    }
    return nil
}
