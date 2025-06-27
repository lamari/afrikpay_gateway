package models

import "testing"

func TestWallet_Validate(t *testing.T) {
    tests := []struct {
        name    string
        wallet  Wallet
        wantErr bool
    }{
        {
            name: "valid wallet",
            wallet: Wallet{
                UserID:   "user123",
                Balance:  100.0,
                Currency: "USDT",
            },
            wantErr: false,
        },
        {
            name: "negative balance",
            wallet: Wallet{
                UserID:   "user123",
                Balance:  -5.0,
                Currency: "USDT",
            },
            wantErr: true,
        },
        {
            name: "unsupported currency",
            wallet: Wallet{
                UserID:   "user123",
                Balance:  10.0,
                Currency: "ETH",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        if err := tt.wallet.Validate(); (err != nil) != tt.wantErr {
            t.Errorf("%s: expected error=%v, got %v", tt.name, tt.wantErr, err)
        }
    }
}
