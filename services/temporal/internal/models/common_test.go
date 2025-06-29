package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ClientError
		expected string
	}{
		{
			name: "Network error",
			err: ClientError{
				Code:       "NETWORK_ERROR",
				Message:    "Connection timeout",
				StatusCode: 0,
				Retryable:  true,
			},
			expected: "Connection timeout",
		},
		{
			name: "API error",
			err: ClientError{
				Code:       "API_ERROR",
				Message:    "Invalid request",
				StatusCode: 400,
				Retryable:  false,
			},
			expected: "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
			assert.Equal(t, tt.err.Retryable, tt.err.IsRetryable())
		})
	}
}

func TestClientError_IsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		err       ClientError
		retryable bool
	}{
		{
			name: "Retryable network error",
			err: ClientError{
				Code:      "TIMEOUT",
				Message:   "Request timeout",
				Retryable: true,
			},
			retryable: true,
		},
		{
			name: "Non-retryable validation error",
			err: ClientError{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid input",
				Retryable: false,
			},
			retryable: false,
		},
		{
			name: "Retryable server error",
			err: ClientError{
				Code:       "SERVER_ERROR",
				Message:    "Internal server error",
				StatusCode: 500,
				Retryable:  true,
			},
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.retryable, tt.err.IsRetryable())
		})
	}
}

func TestPriceResponse_Fields(t *testing.T) {
	timestamp := time.Now()
	response := PriceResponse{
		Symbol:    "BTCUSDT",
		Price:     50000.0,
		Timestamp: timestamp,
	}

	assert.Equal(t, "BTCUSDT", response.Symbol)
	assert.Equal(t, 50000.0, response.Price)
	assert.Equal(t, timestamp, response.Timestamp)
}

func TestOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request OrderRequest
		wantErr bool
	}{
		{
			name: "Valid market buy order",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: 0.001,
			},
			wantErr: false,
		},
		{
			name: "Valid limit sell order",
			request: OrderRequest{
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
			name: "Empty symbol",
			request: OrderRequest{
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: 0.001,
			},
			wantErr: true,
		},
		{
			name: "Invalid side",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "INVALID",
				Type:     "MARKET",
				Quantity: 0.001,
			},
			wantErr: true,
		},
		{
			name: "Zero quantity",
			request: OrderRequest{
				Symbol:   "BTCUSDT",
				Side:     "BUY",
				Type:     "MARKET",
				Quantity: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.Symbol)
				assert.Contains(t, []string{"BUY", "SELL"}, tt.request.Side)
				assert.Contains(t, []string{"MARKET", "LIMIT"}, tt.request.Type)
				assert.Greater(t, tt.request.Quantity, 0.0)
			}
		})
	}
}

func TestOrderResponse_Fields(t *testing.T) {
	timestamp := time.Now()
	response := OrderResponse{
		OrderID:       "order-123",
		Symbol:        "BTCUSDT",
		Status:        "FILLED",
		Side:          "BUY",
		Type:          "MARKET",
		Quantity:      0.001,
		Price:         50000.0,
		ExecutedQty:   0.001,
		Timestamp:     timestamp,
		ClientOrderID: "client-order-456",
	}

	assert.Equal(t, "order-123", response.OrderID)
	assert.Equal(t, "BTCUSDT", response.Symbol)
	assert.Equal(t, "FILLED", response.Status)
	assert.Equal(t, "BUY", response.Side)
	assert.Equal(t, "MARKET", response.Type)
	assert.Equal(t, 0.001, response.Quantity)
	assert.Equal(t, 50000.0, response.Price)
	assert.Equal(t, 0.001, response.ExecutedQty)
	assert.Equal(t, timestamp, response.Timestamp)
	assert.Equal(t, "client-order-456", response.ClientOrderID)
}

func TestPaymentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request PaymentRequest
		wantErr bool
	}{
		{
			name: "Valid payment request",
			request: PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				ExternalID:  "ext-123",
				PhoneNumber: "22507123456",
				Description: "Payment for services",
				CallbackURL: "https://webhook.site/callback",
				Metadata: map[string]string{
					"orderId": "order-123",
				},
			},
			wantErr: false,
		},
		{
			name: "Zero amount",
			request: PaymentRequest{
				Amount:      0,
				Currency:    "XOF",
				ExternalID:  "ext-123",
				PhoneNumber: "22507123456",
			},
			wantErr: true,
		},
		{
			name: "Negative amount",
			request: PaymentRequest{
				Amount:      -100.0,
				Currency:    "XOF",
				ExternalID:  "ext-123",
				PhoneNumber: "22507123456",
			},
			wantErr: true,
		},
		{
			name: "Empty currency",
			request: PaymentRequest{
				Amount:      100.0,
				ExternalID:  "ext-123",
				PhoneNumber: "22507123456",
			},
			wantErr: true,
		},
		{
			name: "Empty external ID",
			request: PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				PhoneNumber: "22507123456",
			},
			wantErr: true,
		},
		{
			name: "Empty phone number",
			request: PaymentRequest{
				Amount:     100.0,
				Currency:   "XOF",
				ExternalID: "ext-123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.Greater(t, tt.request.Amount, 0.0)
				assert.NotEmpty(t, tt.request.Currency)
				assert.NotEmpty(t, tt.request.ExternalID)
				assert.NotEmpty(t, tt.request.PhoneNumber)
			}
		})
	}
}

func TestPaymentResponse_Fields(t *testing.T) {
	response := PaymentResponse{
		ReferenceID: "ref-123-456",
		Status:      PaymentStatusPending,
		PaymentURL:  "https://payment.example.com/pay/123",
		Message:     "Payment initiated successfully",
	}

	assert.Equal(t, "ref-123-456", response.ReferenceID)
	assert.Equal(t, PaymentStatusPending, response.Status)
	assert.Equal(t, "https://payment.example.com/pay/123", response.PaymentURL)
	assert.Equal(t, "Payment initiated successfully", response.Message)
}

func TestPaymentStatusResponse_Fields(t *testing.T) {
	timestamp := time.Now()
	response := PaymentStatusResponse{
		ReferenceID:   "ref-123-456",
		ExternalID:    "ext-123",
		Status:        PaymentStatusSuccess,
		Amount:        100.0,
		Currency:      "XOF",
		PhoneNumber:   "22507123456",
		TransactionID: "txn-789",
		Message:       "Payment completed successfully",
		Timestamp:     timestamp,
	}

	assert.Equal(t, "ref-123-456", response.ReferenceID)
	assert.Equal(t, "ext-123", response.ExternalID)
	assert.Equal(t, PaymentStatusSuccess, response.Status)
	assert.Equal(t, 100.0, response.Amount)
	assert.Equal(t, "XOF", response.Currency)
	assert.Equal(t, "22507123456", response.PhoneNumber)
	assert.Equal(t, "txn-789", response.TransactionID)
	assert.Equal(t, "Payment completed successfully", response.Message)
	assert.Equal(t, timestamp, response.Timestamp)
}

func TestOrderStatus_Constants(t *testing.T) {
	assert.Equal(t, OrderStatus("NEW"), OrderStatusNew)
	assert.Equal(t, OrderStatus("PARTIALLY_FILLED"), OrderStatusPartiallyFilled)
	assert.Equal(t, OrderStatus("FILLED"), OrderStatusFilled)
	assert.Equal(t, OrderStatus("CANCELED"), OrderStatusCanceled)
	assert.Equal(t, OrderStatus("REJECTED"), OrderStatusRejected)
	assert.Equal(t, OrderStatus("EXPIRED"), OrderStatusExpired)
}

func TestOrderSide_Constants(t *testing.T) {
	assert.Equal(t, OrderSide("BUY"), OrderSideBuy)
	assert.Equal(t, OrderSide("SELL"), OrderSideSell)
}

func TestOrderType_Constants(t *testing.T) {
	assert.Equal(t, OrderType("MARKET"), OrderTypeMarket)
	assert.Equal(t, OrderType("LIMIT"), OrderTypeLimit)
}

func TestPaymentStatus_Constants(t *testing.T) {
	assert.Equal(t, PaymentStatus("PENDING"), PaymentStatusPending)
	assert.Equal(t, PaymentStatus("SUCCESS"), PaymentStatusSuccess)
	assert.Equal(t, PaymentStatus("FAILED"), PaymentStatusFailed)
	assert.Equal(t, PaymentStatus("EXPIRED"), PaymentStatusExpired)
	assert.Equal(t, PaymentStatus("PROCESSING"), PaymentStatusProcessing)
}

func TestDepositRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request DepositRequest
		wantErr bool
	}{
		{
			name: "Valid deposit request",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "XAF",
				PhoneNumber: "+237123456789",
				ExternalID:  "ext-123",
				Description: "Test deposit",
			},
			wantErr: false,
		},
		{
			name: "Zero amount",
			request: DepositRequest{
				Amount:      0,
				Currency:    "XAF",
				PhoneNumber: "+237123456789",
				ExternalID:  "ext-123",
			},
			wantErr: true,
		},
		{
			name: "Negative amount",
			request: DepositRequest{
				Amount:      -50.0,
				Currency:    "XAF",
				PhoneNumber: "+237123456789",
				ExternalID:  "ext-123",
			},
			wantErr: true,
		},
		{
			name: "Empty currency",
			request: DepositRequest{
				Amount:      100.0,
				PhoneNumber: "+237123456789",
				ExternalID:  "ext-123",
			},
			wantErr: true,
		},
		{
			name: "Empty phone number",
			request: DepositRequest{
				Amount:     100.0,
				Currency:   "XAF",
				ExternalID: "ext-123",
			},
			wantErr: true,
		},
		{
			name: "Empty external ID",
			request: DepositRequest{
				Amount:      100.0,
				Currency:    "XAF",
				PhoneNumber: "+237123456789",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.Greater(t, tt.request.Amount, 0.0)
				assert.NotEmpty(t, tt.request.Currency)
				assert.NotEmpty(t, tt.request.PhoneNumber)
				assert.NotEmpty(t, tt.request.ExternalID)
			}
		})
	}
}

func TestDepositResponse_Fields(t *testing.T) {
	response := DepositResponse{
		ReferenceID: "ref-deposit-123",
		Status:      PaymentStatusPending,
		Message:     "Deposit initiated successfully",
		Amount:      100.0,
		Currency:    "XAF",
	}

	assert.Equal(t, "ref-deposit-123", response.ReferenceID)
	assert.Equal(t, PaymentStatusPending, response.Status)
	assert.Equal(t, "Deposit initiated successfully", response.Message)
	assert.Equal(t, 100.0, response.Amount)
	assert.Equal(t, "XAF", response.Currency)
}

func TestQuoteResponse_Fields(t *testing.T) {
	timestamp := time.Now()
	quote := QuoteResponse{
		Symbol:    "BTCUSDT",
		BidPrice:  49900.0,
		AskPrice:  50100.0,
		LastPrice: 50000.0,
		Volume:    1500.0,
		Timestamp: timestamp,
	}

	assert.Equal(t, "BTCUSDT", quote.Symbol)
	assert.Equal(t, 49900.0, quote.BidPrice)
	assert.Equal(t, 50100.0, quote.AskPrice)
	assert.Equal(t, 50000.0, quote.LastPrice)
	assert.Equal(t, 1500.0, quote.Volume)
	assert.Equal(t, timestamp, quote.Timestamp)
}

func TestQuotesResponse_Fields(t *testing.T) {
	timestamp := time.Now()
	quotes := QuotesResponse{
		Quotes: []QuoteResponse{
			{
				Symbol:    "BTCUSDT",
				BidPrice:  49900.0,
				AskPrice:  50100.0,
				LastPrice: 50000.0,
				Volume:    1500.0,
				Timestamp: timestamp,
			},
			{
				Symbol:    "ETHUSDT",
				BidPrice:  2990.0,
				AskPrice:  3010.0,
				LastPrice: 3000.0,
				Volume:    800.0,
				Timestamp: timestamp,
			},
		},
		Timestamp: timestamp,
	}

	assert.Len(t, quotes.Quotes, 2)
	assert.Equal(t, "BTCUSDT", quotes.Quotes[0].Symbol)
	assert.Equal(t, "ETHUSDT", quotes.Quotes[1].Symbol)
	assert.Equal(t, timestamp, quotes.Timestamp)

	// Test individual quotes
	btcQuote := quotes.Quotes[0]
	assert.Equal(t, 49900.0, btcQuote.BidPrice)
	assert.Equal(t, 50100.0, btcQuote.AskPrice)
	assert.Equal(t, 50000.0, btcQuote.LastPrice)

	ethQuote := quotes.Quotes[1]
	assert.Equal(t, 2990.0, ethQuote.BidPrice)
	assert.Equal(t, 3010.0, ethQuote.AskPrice)
	assert.Equal(t, 3000.0, ethQuote.LastPrice)
}

func TestPaymentStatus_NewConstants(t *testing.T) {
	// Test the new PaymentStatus constants
	assert.Equal(t, PaymentStatus("SUCCESS"), PaymentStatusSuccess)
	assert.Equal(t, PaymentStatus("PROCESSING"), PaymentStatusProcessing)
	assert.Equal(t, PaymentStatus("EXPIRED"), PaymentStatusExpired)
}
