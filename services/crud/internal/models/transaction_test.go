package models

import "testing"

func TestTransaction_Validate(t *testing.T) {
    tests := []struct {
        name        string
        transaction Transaction
        wantErr     bool
    }{
        {
            name: "valid deposit",
            transaction: Transaction{
                WalletID: "wallet1",
                Amount:   50.0,
                Type:     TransactionDeposit,
                Status:   StatusPending,
            },
            wantErr: false,
        },
        {
            name: "zero amount",
            transaction: Transaction{
                WalletID: "wallet1",
                Amount:   0,
                Type:     TransactionDeposit,
                Status:   StatusPending,
            },
            wantErr: true,
        },
        {
            name: "invalid type",
            transaction: Transaction{
                WalletID: "wallet1",
                Amount:   10,
                Type:     "stake",
                Status:   StatusPending,
            },
            wantErr: true,
        },
        {
            name: "invalid status",
            transaction: Transaction{
                WalletID: "wallet1",
                Amount:   10,
                Type:     TransactionDeposit,
                Status:   "unknown",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        if err := tt.transaction.Validate(); (err != nil) != tt.wantErr {
            t.Errorf("%s: expected error=%v, got %v", tt.name, tt.wantErr, err)
        }
    }
}
