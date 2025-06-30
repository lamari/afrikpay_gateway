package activities

import (
	"context"

	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/models"
)

// BinanceActivities wraps Binance client functionality for Temporal activities
type BinanceActivities struct {
	client clients.BinanceClient
}

// NewBinanceActivities creates a new BinanceActivities instance
func NewBinanceActivities(client clients.BinanceClient) *BinanceActivities {
	return &BinanceActivities{
		client: client,
	}
}

// GetPrice gets the current price for a symbol
func (a *BinanceActivities) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	return a.client.GetPrice(ctx, symbol)
}

// PlaceOrder places a new order
func (a *BinanceActivities) PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	return a.client.PlaceOrder(ctx, request)
}

// GetOrderStatus gets the status of an order
func (a *BinanceActivities) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error) {
	return a.client.GetOrderStatus(ctx, symbol, orderID)
}

// GetQuotes gets market quotes for all symbols
func (a *BinanceActivities) GetQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	return a.client.GetQuotes(ctx)
}

// GetQuote gets market quote for a specific symbol
func (a *BinanceActivities) GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error) {
	return a.client.GetQuote(ctx, symbol)
}

// GetAllOrders gets all open orders for the account
func (a *BinanceActivities) GetAllOrders(ctx context.Context) (*models.OrdersResponse, error) {
	return a.client.GetAllOrders(ctx)
}

// HealthCheck performs a health check
func (a *BinanceActivities) HealthCheck(ctx context.Context) error {
	return a.client.HealthCheck(ctx)
}
