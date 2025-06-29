package workflows

import (
	"testing"
	"time"

	"github.com/afrikpay/gateway/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func TestBinancePriceWorkflow_Success(t *testing.T) {
	// Given
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock the activity
	expectedPrice := &models.PriceResponse{
		Symbol:    "BTCUSDT",
		Price:     45000.50,
		Timestamp: time.Now(),
		Success:   true,
	}

	env.OnActivity("GetBinancePrice", mock.Anything, "BTCUSDT").Return(expectedPrice, nil)

	// When
	env.ExecuteWorkflow(BinancePriceWorkflow, "BTCUSDT")

	// Then
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())

	var result models.PriceResponse
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, 45000.50, result.Price)
	assert.True(t, result.Success)
}

func TestBinancePriceWorkflow_ActivityError(t *testing.T) {
	// Given
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock the activity to return an error
	env.OnActivity("GetBinancePrice", mock.Anything, "INVALID").Return(nil, assert.AnError)

	// When
	env.ExecuteWorkflow(BinancePriceWorkflow, "INVALID")

	// Then
	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())
}

func TestBinancePriceWorkflow_EmptySymbol(t *testing.T) {
	// Given
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// When
	env.ExecuteWorkflow(BinancePriceWorkflow, "")

	// Then
	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())
}
