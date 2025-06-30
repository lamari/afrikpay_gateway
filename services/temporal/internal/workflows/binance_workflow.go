package workflows

import (
	"fmt"
	"time"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/models"
	"go.temporal.io/sdk/workflow"
)

// BinancePriceWorkflowInput represents the input for the Binance price workflow
type BinancePriceWorkflowInput struct {
	Symbol string `json:"symbol"`
}

// BinancePriceWorkflow is a simple workflow that gets the price for a symbol from Binance
func BinancePriceWorkflow(ctx workflow.Context, symbol string) (*models.PriceResponse, error) {
	// Validate input
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	// Get workflow logger
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Binance price workflow", "symbol", symbol)

	// Set up activity options
	activityOptions := defaultActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute GetPrice activity
	var priceResponse *models.PriceResponse

	err := workflow.ExecuteActivity(ctx, activities.GetBinanceActivitiesFromFactory().GetPrice, symbol).Get(ctx, &priceResponse)
	if err != nil {
		logger.Error("Failed to get price from Binance", "symbol", symbol, "error", err)
		return nil, fmt.Errorf("failed to get price for symbol %s: %w", symbol, err)
	}

	logger.Info("Successfully retrieved price from Binance",
		"symbol", priceResponse.Symbol,
		"price", priceResponse.Price,
		"success", priceResponse.Success)

	return priceResponse, nil
}

// BinanceQuotesWorkflow retrieves multiple quotes from Binance
func BinanceQuotesWorkflow(ctx workflow.Context) (*models.QuotesResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("BinanceQuotesWorkflow started")

	// Configure activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute GetQuotes activity
	var quotesResponse models.QuotesResponse
	err := workflow.ExecuteActivity(ctx, "GetQuotes").Get(ctx, &quotesResponse)
	if err != nil {
		logger.Error("Failed to get Binance quotes", "error", err)
		return nil, err
	}

	logger.Info("BinanceQuotesWorkflow completed successfully", "quotes_count", len(quotesResponse.Quotes))
	return &quotesResponse, nil
}

// BinanceOrdersWorkflow retrieves all orders from Binance
func BinanceOrdersWorkflow(ctx workflow.Context) (*models.OrdersResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("BinanceOrdersWorkflow started")

	// Configure activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute GetAllOrders activity
	var ordersResponse models.OrdersResponse
	err := workflow.ExecuteActivity(ctx, "GetAllOrders").Get(ctx, &ordersResponse)
	if err != nil {
		logger.Error("Failed to get Binance orders", "error", err)
		return nil, err
	}

	logger.Info("BinanceOrdersWorkflow completed successfully", "orders_count", len(ordersResponse.Orders))
	return &ordersResponse, nil
}

// BinancePriceWorkflowWithInput is a wrapper that accepts structured input
func BinancePriceWorkflowWithInput(ctx workflow.Context, input BinancePriceWorkflowInput) (*models.PriceResponse, error) {
	return BinancePriceWorkflow(ctx, input.Symbol)
}
