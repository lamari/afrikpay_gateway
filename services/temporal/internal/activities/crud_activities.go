package activities

import (
	"context"

	"github.com/afrikpay/gateway/internal/models"
)

// CrudActivities wraps CRUD client functionality for Temporal activities
type CrudActivities struct {
	client models.CrudClient
}

// NewCrudActivities creates a new CrudActivities instance
func NewCrudActivities(client models.CrudClient) *CrudActivities {
	return &CrudActivities{
		client: client,
	}
}

// CreateTransaction creates a new transaction in the CRUD service
func (a *CrudActivities) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.TransactionResponse, error) {
	return a.client.CreateTransaction(ctx, transaction)
}

// UpdateWalletBalance updates a wallet's balance in the CRUD service
func (a *CrudActivities) UpdateWalletBalance(ctx context.Context, walletID string, amount float64, currency string) (*models.WalletResponse, error) {
	return a.client.UpdateWalletBalance(ctx, walletID, amount, currency)
}

// GetWallet retrieves wallet information by user ID and currency
func (a *CrudActivities) GetWallet(ctx context.Context, userID string, currency string) (*models.WalletResponse, error) {
	return a.client.GetWallet(ctx, userID, currency)
}

// HealthCheck performs a health check against the CRUD service
func (a *CrudActivities) HealthCheck(ctx context.Context) error {
	return a.client.HealthCheck(ctx)
}
