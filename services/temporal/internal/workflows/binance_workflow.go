package workflows

import (
	"fmt"

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

// BinancePriceWorkflowWithInput is a wrapper that accepts structured input
func BinancePriceWorkflowWithInput(ctx workflow.Context, input BinancePriceWorkflowInput) (*models.PriceResponse, error) {
	return BinancePriceWorkflow(ctx, input.Symbol)
}
