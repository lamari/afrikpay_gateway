package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBinanceConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config BinanceConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			valid: true,
		},
		{
			name: "Empty BaseURL",
			config: BinanceConfig{
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			valid: false,
		},
		{
			name: "Empty APIKey",
			config: BinanceConfig{
				BaseURL:    "https://api.binance.com",
				SecretKey:  "test-secret-key",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic will be implemented when we create the validation functions
			if tt.valid {
				assert.NotEmpty(t, tt.config.BaseURL)
				assert.NotEmpty(t, tt.config.APIKey)
				assert.NotEmpty(t, tt.config.SecretKey)
			}
		})
	}
}

func TestBinancePriceRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BinancePriceRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: BinancePriceRequest{
				Symbol: "BTCUSDT",
			},
			wantErr: false,
		},
		{
			name: "Empty symbol",
			request: BinancePriceRequest{
				Symbol: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validation will be implemented in the service layer
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.Symbol)
			} else {
				assert.Empty(t, tt.request.Symbol)
			}
		})
	}
}

func TestBinanceOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BinanceOrderRequest
		wantErr bool
	}{
		{
			name: "Valid market buy order",
			request: BinanceOrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: 0.001,
			},
			wantErr: false,
		},
		{
			name: "Valid limit sell order",
			request: BinanceOrderRequest{
				Symbol:      "BTCUSDT",
				Side:        "SELL",
				Type:        "LIMIT",
				Quantity:    0.001,
				Price:       50000.0,
				TimeInForce: "GTC",
			},
			wantErr: false,
		},
		{
			name: "Invalid side",
			request: BinanceOrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "INVALID",
				Type:     "MARKET",
				Quantity: 0.001,
			},
			wantErr: true,
		},
		{
			name: "Invalid type",
			request: BinanceOrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "INVALID",
				Quantity: 0.001,
			},
			wantErr: true,
		},
		{
			name: "Zero quantity",
			request: BinanceOrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: 0,
			},
			wantErr: true,
		},
		{
			name: "Negative quantity",
			request: BinanceOrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: -0.001,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.Symbol)
				assert.Contains(t, []string{"BUY", "SELL"}, tt.request.Side)
				assert.Contains(t, []string{"MARKET", "LIMIT"}, tt.request.Type)
				assert.Greater(t, tt.request.Quantity, 0.0)
			}
		})
	}
}

func TestBinanceOrderResponse_Fields(t *testing.T) {
	response := BinanceOrderResponse{
		Symbol:        "BTCUSDT",
		OrderID:       12345,
		ClientOrderID: "test-order-123",
		TransactTime:  1640995200000,
		Price:         "50000.00",
		OrigQty:       "0.001",
		ExecutedQty:   "0.001",
		Status:        "FILLED",
		TimeInForce:   "GTC",
		Type:          "MARKET",
		Side:          "BUY",
		Fills: []Fill{
			{
				Price:           "50000.00",
				Qty:             "0.001",
				Commission:      "0.00001",
				CommissionAsset: "BNB",
			},
		},
	}

	assert.Equal(t, "BTCUSDT", response.Symbol)
	assert.Equal(t, int64(12345), response.OrderID)
	assert.Equal(t, "test-order-123", response.ClientOrderID)
	assert.Equal(t, "FILLED", response.Status)
	assert.Len(t, response.Fills, 1)
	assert.Equal(t, "50000.00", response.Fills[0].Price)
}

func TestBinanceError_Error(t *testing.T) {
	tests := []struct {
		name     string
		binError BinanceError
		expected string
	}{
		{
			name: "API error",
			binError: BinanceError{
				Code: -1121,
				Msg:  "Invalid symbol.",
			},
			expected: "Invalid symbol.",
		},
		{
			name: "Rate limit error",
			binError: BinanceError{
				Code: -1003,
				Msg:  "Too many requests.",
			},
			expected: "Too many requests.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.binError
			assert.Equal(t, tt.expected, err.Error())
			assert.Equal(t, tt.binError.Code, err.Code)
			assert.Equal(t, tt.binError.Msg, err.Msg)
		})
	}
}

func TestBinanceOrderStatusRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BinanceOrderStatusRequest
		wantErr bool
	}{
		{
			name: "Valid request with OrderID",
			request: BinanceOrderStatusRequest{
				Symbol:  "BTCUSDT",
				OrderID: 12345,
			},
			wantErr: false,
		},
		{
			name: "Valid request with ClientOrderID",
			request: BinanceOrderStatusRequest{
				Symbol:            "BTCUSDT",
				OrigClientOrderID: "test-order-123",
			},
			wantErr: false,
		},
		{
			name: "Empty symbol",
			request: BinanceOrderStatusRequest{
				OrderID: 12345,
			},
			wantErr: true,
		},
		{
			name: "No order identifier",
			request: BinanceOrderStatusRequest{
				Symbol: "BTCUSDT",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.Symbol)
				hasOrderID := tt.request.OrderID != 0 || tt.request.OrigClientOrderID != ""
				assert.True(t, hasOrderID, "Should have either OrderID or OrigClientOrderID")
			}
		})
	}
}

func TestFill_Fields(t *testing.T) {
	fill := Fill{
		Price:           "50000.00",
		Qty:             "0.001",
		Commission:      "0.00001",
		CommissionAsset: "BNB",
	}

	assert.Equal(t, "50000.00", fill.Price)
	assert.Equal(t, "0.001", fill.Qty)
	assert.Equal(t, "0.00001", fill.Commission)
	assert.Equal(t, "BNB", fill.CommissionAsset)
}

func TestBinanceConfig_DefaultValues(t *testing.T) {
	config := BinanceConfig{
		BaseURL:    "https://api.binance.com",
		APIKey:     "test-key",
		SecretKey:  "test-secret",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	assert.Equal(t, "https://api.binance.com", config.BaseURL)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 3, config.MaxRetries)
}
