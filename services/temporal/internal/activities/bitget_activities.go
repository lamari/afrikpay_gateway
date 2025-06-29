package activities

import (
	"context"

	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/models"
)

// BitgetActivities wraps Bitget client functionality for Temporal activities
type BitgetActivities struct {
	client clients.BitgetClient
}

// NewBitgetActivities creates a new BitgetActivities instance
func NewBitgetActivities(client clients.BitgetClient) *BitgetActivities {
	return &BitgetActivities{
		client: client,
	}
}

// GetPrice gets the current price for a symbol
func (a *BitgetActivities) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	return a.client.GetPrice(ctx, symbol)
}

// PlaceOrder places a new order
func (a *BitgetActivities) PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	return a.client.PlaceOrder(ctx, request)
}

// GetOrderStatus gets the status of an order
func (a *BitgetActivities) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error) {
	return a.client.GetOrderStatus(ctx, symbol, orderID)
}

// GetQuotes gets market quotes for all symbols
func (a *BitgetActivities) GetQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	return a.client.GetQuotes(ctx)
}

// GetQuote gets market quote for a specific symbol
func (a *BitgetActivities) GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error) {
	return a.client.GetQuote(ctx, symbol)
}

// HealthCheck performs a health check
func (a *BitgetActivities) HealthCheck(ctx context.Context) error {
	// BitgetClient étend ClientInterface qui contient HealthCheck
	return a.client.HealthCheck(ctx)
}
