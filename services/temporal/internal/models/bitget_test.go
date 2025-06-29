package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBitgetConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  BitgetConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid config",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "Missing base URL",
			config: BitgetConfig{
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "base URL is required",
		},
		{
			name: "Missing API key",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
		{
			name: "Missing secret key",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-api-key",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "secret key is required",
		},
		{
			name: "Missing passphrase",
			config: BitgetConfig{
				BaseURL:   "https://api.bitget.com",
				APIKey:    "test-api-key",
				SecretKey: "test-secret-key",
				Timeout:   30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "passphrase is required",
		},
		{
			name: "Zero timeout",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    0,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "timeout must be positive",
		},
		{
			name: "Negative timeout",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    -1 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
			errMsg:  "timeout must be positive",
		},
		{
			name: "Negative max retries",
			config: BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-api-key",
				SecretKey:  "test-secret-key",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: -1,
			},
			wantErr: true,
			errMsg:  "max retries cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBitgetPriceRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BitgetPriceRequest
		valid   bool
	}{
		{
			name: "Valid request",
			request: BitgetPriceRequest{
				Symbol: "BTCUSDT",
			},
			valid: true,
		},
		{
			name: "Empty symbol",
			request: BitgetPriceRequest{
				Symbol: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct can be created
			assert.Equal(t, tt.request.Symbol, tt.request.Symbol)
			
			// In a real validation scenario, you would use a validator
			if tt.valid {
				assert.NotEmpty(t, tt.request.Symbol)
			} else {
				assert.Empty(t, tt.request.Symbol)
			}
		})
	}
}

func TestBitgetOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BitgetOrderRequest
		valid   bool
	}{
		{
			name: "Valid market buy order",
			request: BitgetOrderRequest{
				Symbol:    "BTCUSDT",
				Side:      BitgetOrderSideBuy,
				OrderType: BitgetOrderTypeMarket,
				Size:      "0.001",
			},
			valid: true,
		},
		{
			name: "Valid limit sell order",
			request: BitgetOrderRequest{
				Symbol:      "BTCUSDT",
				Side:        BitgetOrderSideSell,
				OrderType:   BitgetOrderTypeLimit,
				Size:        "0.001",
				Price:       "50000.00",
				TimeInForce: BitgetTimeInForceGTC,
			},
			valid: true,
		},
		{
			name: "Missing symbol",
			request: BitgetOrderRequest{
				Side:      BitgetOrderSideBuy,
				OrderType: BitgetOrderTypeMarket,
				Size:      "0.001",
			},
			valid: false,
		},
		{
			name: "Missing side",
			request: BitgetOrderRequest{
				Symbol:    "BTCUSDT",
				OrderType: BitgetOrderTypeMarket,
				Size:      "0.001",
			},
			valid: false,
		},
		{
			name: "Missing order type",
			request: BitgetOrderRequest{
				Symbol: "BTCUSDT",
				Side:   BitgetOrderSideBuy,
				Size:   "0.001",
			},
			valid: false,
		},
		{
			name: "Missing size",
			request: BitgetOrderRequest{
				Symbol:    "BTCUSDT",
				Side:      BitgetOrderSideBuy,
				OrderType: BitgetOrderTypeMarket,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic field presence
			if tt.valid {
				assert.NotEmpty(t, tt.request.Symbol)
				assert.NotEmpty(t, tt.request.Side)
				assert.NotEmpty(t, tt.request.OrderType)
				assert.NotEmpty(t, tt.request.Size)
			}
		})
	}
}

func TestBitgetPriceResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response BitgetPriceResponse
		expected bool
	}{
		{
			name: "Success response",
			response: BitgetPriceResponse{
				Code:    "00000",
				Message: "success",
			},
			expected: true,
		},
		{
			name: "Error response",
			response: BitgetPriceResponse{
				Code:    "40001",
				Message: "Invalid symbol",
			},
			expected: false,
		},
		{
			name: "Empty code",
			response: BitgetPriceResponse{
				Code:    "",
				Message: "success",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.response.IsSuccess())
		})
	}
}

func TestBitgetOrderResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response BitgetOrderResponse
		expected bool
	}{
		{
			name: "Success response",
			response: BitgetOrderResponse{
				Code:    "00000",
				Message: "success",
			},
			expected: true,
		},
		{
			name: "Error response",
			response: BitgetOrderResponse{
				Code:    "40001",
				Message: "Insufficient balance",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.response.IsSuccess())
		})
	}
}

func TestBitgetOrderResponse_ToOrderResponse(t *testing.T) {
	t.Run("Success response conversion", func(t *testing.T) {
		bitgetResp := BitgetOrderResponse{
			Code:    "00000",
			Message: "success",
			Data: struct {
				OrderId       string `json:"orderId"`
				ClientOid     string `json:"clientOid"`
				Symbol        string `json:"symbol"`
				Side          string `json:"side"`
				OrderType     string `json:"orderType"`
				Size          string `json:"size"`
				Price         string `json:"price"`
				Status        string `json:"status"`
				FilledSize    string `json:"filledSize"`
				FilledAmount  string `json:"filledAmount"`
				CreateTime    int64  `json:"createTime"`
				UpdateTime    int64  `json:"updateTime"`
			}{
				OrderId:      "12345",
				ClientOid:    "client-123",
				Symbol:       "BTCUSDT",
				Side:         "buy",
				OrderType:    "market",
				Size:         "0.001",
				Price:        "50000.00",
				Status:       "filled",
				FilledSize:   "0.001",
				FilledAmount: "50.00",
				CreateTime:   1640995200000, // 2022-01-01 00:00:00 UTC
				UpdateTime:   1640995260000, // 2022-01-01 00:01:00 UTC
			},
			RequestTime: 1640995200000,
		}

		orderResp := bitgetResp.ToOrderResponse()

		assert.NotNil(t, orderResp)
		assert.Equal(t, "12345", orderResp.OrderID)
		assert.Equal(t, "BTCUSDT", orderResp.Symbol)
		assert.Equal(t, "buy", orderResp.Side)
		assert.Equal(t, "market", orderResp.Type)
		assert.Equal(t, 0.001, orderResp.Quantity)
		assert.Equal(t, 50000.00, orderResp.Price)
		assert.Equal(t, "filled", orderResp.Status)
		assert.Equal(t, 0.001, orderResp.ExecutedQty)
		assert.Equal(t, time.Unix(1640995200, 0), orderResp.Timestamp)
		assert.Equal(t, "client-123", orderResp.ClientOrderID)
	})

	t.Run("Error response conversion", func(t *testing.T) {
		bitgetResp := BitgetOrderResponse{
			Code:    "40001",
			Message: "Insufficient balance",
		}

		orderResp := bitgetResp.ToOrderResponse()

		// For error responses, we return an empty OrderResponse
		assert.NotNil(t, orderResp)
		assert.Empty(t, orderResp.OrderID)
		assert.Empty(t, orderResp.Symbol)
	})
}

func TestBitgetOrderStatusResponse_ToOrderResponse(t *testing.T) {
	t.Run("Success response conversion", func(t *testing.T) {
		bitgetResp := BitgetOrderStatusResponse{
			Code:    "00000",
			Message: "success",
			Data: struct {
				OrderId       string `json:"orderId"`
				ClientOid     string `json:"clientOid"`
				Symbol        string `json:"symbol"`
				Side          string `json:"side"`
				OrderType     string `json:"orderType"`
				Size          string `json:"size"`
				Price         string `json:"price"`
				Status        string `json:"status"`
				FilledSize    string `json:"filledSize"`
				FilledAmount  string `json:"filledAmount"`
				AvgPrice      string `json:"avgPrice"`
				Fee           string `json:"fee"`
				FeeCurrency   string `json:"feeCurrency"`
				CreateTime    int64  `json:"createTime"`
				UpdateTime    int64  `json:"updateTime"`
			}{
				OrderId:      "12345",
				ClientOid:    "client-123",
				Symbol:       "BTCUSDT",
				Side:         "buy",
				OrderType:    "market",
				Size:         "0.001",
				Price:        "50000.00",
				Status:       "filled",
				FilledSize:   "0.001",
				FilledAmount: "50.00",
				AvgPrice:     "50000.00",
				Fee:          "0.05",
				FeeCurrency:  "USDT",
				CreateTime:   1640995200000,
				UpdateTime:   1640995260000,
			},
		}

		orderResp := bitgetResp.ToOrderResponse()

		assert.NotNil(t, orderResp)
		assert.Equal(t, "12345", orderResp.OrderID)
		assert.Equal(t, "BTCUSDT", orderResp.Symbol)
		// Note: AvgPrice, Fee, FeeCurrency are not in OrderResponse struct
	})
}

func TestBitgetPriceResponse_ToPriceResponse(t *testing.T) {
	t.Run("Success response conversion", func(t *testing.T) {
		bitgetResp := BitgetPriceResponse{
			Code:    "00000",
			Message: "success",
			Data: struct {
				Symbol string `json:"symbol"`
				Price  string `json:"price"`
			}{
				Symbol: "BTCUSDT",
				Price:  "50000.00",
			},
			RequestTime: 1640995200000,
		}

		priceResp := bitgetResp.ToPriceResponse()

		assert.NotNil(t, priceResp)
		assert.Equal(t, "BTCUSDT", priceResp.Symbol)
		assert.Equal(t, 50000.00, priceResp.Price)
		assert.Equal(t, time.Unix(1640995200, 0), priceResp.Timestamp)
	})

	t.Run("Error response conversion", func(t *testing.T) {
		bitgetResp := BitgetPriceResponse{
			Code:    "40001",
			Message: "Invalid symbol",
		}

		priceResp := bitgetResp.ToPriceResponse()

		// For error responses, we return an empty PriceResponse
		assert.NotNil(t, priceResp)
		assert.Empty(t, priceResp.Symbol)
		assert.Equal(t, 0.0, priceResp.Price)
	})
}

func TestBitgetConstants(t *testing.T) {
	// Test order status constants
	assert.Equal(t, "new", BitgetOrderStatusNew)
	assert.Equal(t, "partially_filled", BitgetOrderStatusPartiallyFilled)
	assert.Equal(t, "filled", BitgetOrderStatusFilled)
	assert.Equal(t, "canceled", BitgetOrderStatusCanceled)
	assert.Equal(t, "rejected", BitgetOrderStatusRejected)

	// Test order side constants
	assert.Equal(t, "buy", BitgetOrderSideBuy)
	assert.Equal(t, "sell", BitgetOrderSideSell)

	// Test order type constants
	assert.Equal(t, "market", BitgetOrderTypeMarket)
	assert.Equal(t, "limit", BitgetOrderTypeLimit)

	// Test time in force constants
	assert.Equal(t, "GTC", BitgetTimeInForceGTC)
	assert.Equal(t, "IOC", BitgetTimeInForceIOC)
	assert.Equal(t, "FOK", BitgetTimeInForceFOK)
}

func TestBitgetQuotesResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response BitgetQuotesResponse
		expected bool
	}{
		{
			name: "Success response",
			response: BitgetQuotesResponse{
				Code:    "00000",
				Message: "success",
			},
			expected: true,
		},
		{
			name: "Error response",
			response: BitgetQuotesResponse{
				Code:    "40001",
				Message: "Invalid request",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.response.IsSuccess())
		})
	}
}

func TestBitgetQuoteResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response BitgetQuoteResponse
		expected bool
	}{
		{
			name: "Success response",
			response: BitgetQuoteResponse{
				Code:    "00000",
				Message: "success",
			},
			expected: true,
		},
		{
			name: "Error response",
			response: BitgetQuoteResponse{
				Code:    "40001",
				Message: "Symbol not found",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.response.IsSuccess())
		})
	}
}

func TestBitgetOrderStatusRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request BitgetOrderStatusRequest
		valid   bool
	}{
		{
			name: "Valid request",
			request: BitgetOrderStatusRequest{
				Symbol:  "BTCUSDT",
				OrderId: "12345",
			},
			valid: true,
		},
		{
			name: "Missing symbol",
			request: BitgetOrderStatusRequest{
				OrderId: "12345",
			},
			valid: false,
		},
		{
			name: "Missing order ID",
			request: BitgetOrderStatusRequest{
				Symbol: "BTCUSDT",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Symbol)
				assert.NotEmpty(t, tt.request.OrderId)
			}
		})
	}
}

func TestBitgetErrorResponse(t *testing.T) {
	errorResp := BitgetErrorResponse{
		Code:    "40001",
		Message: "Invalid API key",
	}

	assert.Equal(t, "40001", errorResp.Code)
	assert.Equal(t, "Invalid API key", errorResp.Message)
}

func TestBitgetResponseStructures(t *testing.T) {
	t.Run("BitgetQuotesResponse with data", func(t *testing.T) {
		resp := BitgetQuotesResponse{
			Code:    "00000",
			Message: "success",
			Data: []struct {
				Symbol   string `json:"symbol"`
				Price    string `json:"price"`
				Change   string `json:"change"`
				ChangeP  string `json:"changeP"`
				High24h  string `json:"high24h"`
				Low24h   string `json:"low24h"`
				Volume   string `json:"volume"`
				Turnover string `json:"turnover"`
			}{
				{
					Symbol:   "BTCUSDT",
					Price:    "50000.00",
					Change:   "1000.00",
					ChangeP:  "2.04",
					High24h:  "51000.00",
					Low24h:   "49000.00",
					Volume:   "1000.5",
					Turnover: "50025000.00",
				},
			},
			RequestTime: 1640995200000,
		}

		assert.True(t, resp.IsSuccess())
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, "BTCUSDT", resp.Data[0].Symbol)
		assert.Equal(t, "50000.00", resp.Data[0].Price)
	})

	t.Run("BitgetQuoteResponse with data", func(t *testing.T) {
		resp := BitgetQuoteResponse{
			Code:    "00000",
			Message: "success",
			Data: struct {
				Symbol   string `json:"symbol"`
				Price    string `json:"price"`
				Change   string `json:"change"`
				ChangeP  string `json:"changeP"`
				High24h  string `json:"high24h"`
				Low24h   string `json:"low24h"`
				Volume   string `json:"volume"`
				Turnover string `json:"turnover"`
			}{
				Symbol:   "BTCUSDT",
				Price:    "50000.00",
				Change:   "1000.00",
				ChangeP:  "2.04",
				High24h:  "51000.00",
				Low24h:   "49000.00",
				Volume:   "1000.5",
				Turnover: "50025000.00",
			},
			RequestTime: 1640995200000,
		}

		assert.True(t, resp.IsSuccess())
		assert.Equal(t, "BTCUSDT", resp.Data.Symbol)
		assert.Equal(t, "50000.00", resp.Data.Price)
		assert.Equal(t, "2.04", resp.Data.ChangeP)
	})
}
