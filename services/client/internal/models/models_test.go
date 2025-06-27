package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDepositRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request DepositRequest
		wantErr bool
	}{
		{
			name: "Valid deposit request",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "XOF",
				PhoneNumber: "+22670000000",
				ExternalID:  "test-ext-123",
				Description: "Test payment",
			},
			wantErr: false,
		},
		{
			name: "Invalid amount - zero",
			request: DepositRequest{
				Amount:      0,
				Currency:    "XOF",
				PhoneNumber: "+22670000000",
				Description: "Test payment",
			},
			wantErr: true,
		},
		{
			name: "Invalid amount - negative",
			request: DepositRequest{
				Amount:      -50.0,
				Currency:    "XOF",
				PhoneNumber: "+22670000000",
				Description: "Test payment",
			},
			wantErr: true,
		},
		{
			name: "Invalid currency - empty",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "",
				PhoneNumber: "+22670000000",
				Description: "Test payment",
			},
			wantErr: true,
		},
		{
			name: "Invalid phone number - empty",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "XOF",
				PhoneNumber: "",
				Description: "Test payment",
			},
			wantErr: true,
		},
		{
			name: "Invalid phone number - wrong format",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "XOF",
				PhoneNumber: "70000000",
				Description: "Test payment",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request OrderRequest
		wantErr bool
	}{
		{
			name: "Valid buy limit order",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: 0.1,
				Price:    41000.0,
			},
			wantErr: false,
		},
		{
			name: "Valid sell market order",
			request: OrderRequest{
				Symbol:   "ETHUSDT",
				Side:     "SELL",
				Type:     "MARKET",
				Quantity: 1.0,
			},
			wantErr: false,
		},
		{
			name: "Invalid symbol - empty",
			request: OrderRequest{
				Symbol:   "",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: 0.1,
				Price:    41000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid side",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "invalid",
				Type:     "LIMIT",
				Quantity: 0.1,
				Price:    41000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid type",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "invalid",
				Quantity: 0.1,
				Price:    41000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid quantity - zero",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: 0,
				Price:    41000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid quantity - negative",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: -0.1,
				Price:    41000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid price for limit order - zero",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: 0.1,
				Price:    0,
			},
			wantErr: true,
		},
		{
			name: "Invalid price for limit order - negative",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "LIMIT",
				Quantity: 0.1,
				Price:    -1000.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaymentStatus_String(t *testing.T) {
	tests := []struct {
		status   PaymentStatus
		expected string
	}{
		{PaymentStatusPending, "PENDING"},
		{PaymentStatusSuccess, "SUCCESS"},
		{PaymentStatusFailed, "FAILED"},
		{PaymentStatusCanceled, "CANCELED"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}



func TestOrderStatus_String(t *testing.T) {
	tests := []struct {
		status OrderStatus
		want   string
	}{
		{OrderStatusNew, "NEW"},
		{OrderStatusPartiallyFilled, "PARTIALLY_FILLED"},
		{OrderStatusFilled, "FILLED"},
		{OrderStatusCanceled, "CANCELED"},
		{OrderStatusRejected, "REJECTED"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}



func TestResilienceStats_Reset(t *testing.T) {
	stats := &ResilienceStats{
		TotalRequests:      100,
		SuccessfulRequests: 80,
		FailedRequests:     20,
		LastReset:          time.Now().Add(-time.Hour),
	}

	oldLastReset := stats.LastReset
	stats.Reset()

	assert.Equal(t, int64(0), stats.TotalRequests)
	assert.Equal(t, int64(0), stats.SuccessfulRequests)
	assert.Equal(t, int64(0), stats.FailedRequests)
	assert.True(t, stats.LastReset.After(oldLastReset))
}

func TestMobileMoneyConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  MobileMoneyConfig
		wantErr bool
	}{
		{
			name: "Valid MTN config",
			config: MobileMoneyConfig{
				Provider:        "mtn",
				BaseURL:         "https://sandbox.momodeveloper.mtn.com",
				APIKey:          "test-key",
				APISecret:       "test-secret",
				SubscriptionKey: "test-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			wantErr: false,
		},
		{
			name: "Valid Orange config",
			config: MobileMoneyConfig{
				Provider:        "orange",
				BaseURL:         "https://api.orange.com",
				APIKey:          "test-key",
				APISecret:       "test-secret",
				SubscriptionKey: "test-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			wantErr: false,
		},
		{
			name: "Invalid provider",
			config: MobileMoneyConfig{
				Provider:        "invalid",
				BaseURL:         "https://example.com",
				APIKey:          "test-key",
				APISecret:       "test-secret",
				SubscriptionKey: "test-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			wantErr: true,
		},
		{
			name: "Empty base URL",
			config: MobileMoneyConfig{
				Provider:        "mtn",
				BaseURL:         "",
				APIKey:          "test-key",
				APISecret:       "test-secret",
				SubscriptionKey: "test-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			wantErr: true,
		},
		{
			name: "Empty API key",
			config: MobileMoneyConfig{
				Provider:        "mtn",
				BaseURL:         "https://sandbox.momodeveloper.mtn.com",
				APIKey:          "",
				APISecret:       "test-secret",
				SubscriptionKey: "test-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBinanceConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  BinanceConfig
		wantErr bool
	}{
		{
			name: "Valid config",
			config: BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "Empty base URL",
			config: BinanceConfig{
				BaseURL:    "",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
		},
		{
			name: "Empty API key",
			config: BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
		},
		{
			name: "Empty secret key",
			config: BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-key",
				SecretKey:  "",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBitgetConfig_Validate is now in bitget_test.go
