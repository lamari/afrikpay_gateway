package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMobileMoneyConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config MobileMoneyConfig
		valid  bool
	}{
		{
			name: "Valid MTN config",
			config: MobileMoneyConfig{
				BaseURL:         "https://sandbox.momodeveloper.mtn.com",
				APIKey:          "test-api-key",
				APISecret:       "test-api-secret",
				SubscriptionKey: "test-subscription-key",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			},
			valid: true,
		},
		{
			name: "Valid Orange config",
			config: MobileMoneyConfig{
				BaseURL:    "https://api.orange.com",
				APIKey:     "test-merchant-key",
				APISecret:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			valid: true,
		},
		{
			name: "Empty BaseURL",
			config: MobileMoneyConfig{
				APIKey:     "test-api-key",
				APISecret:  "test-api-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.config.BaseURL)
				assert.NotEmpty(t, tt.config.APIKey)
			}
		})
	}
}

func TestMTNPaymentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request MTNPaymentRequest
		wantErr bool
	}{
		{
			name: "Valid MTN payment request",
			request: MTNPaymentRequest{
				Amount:     "100",
				Currency:   "EUR",
				ExternalID: "test-external-id-123",
				Payer: MTNPayer{
					PartyIDType: "MSISDN",
					PartyID:     "46733123450",
				},
				PayerMessage: "Payment for services",
				PayeeNote:    "Thank you for your payment",
				CallbackURL:  "https://webhook.site/callback",
			},
			wantErr: false,
		},
		{
			name: "Empty amount",
			request: MTNPaymentRequest{
				Currency:   "EUR",
				ExternalID: "test-external-id-123",
				Payer: MTNPayer{
					PartyIDType: "MSISDN",
					PartyID:     "46733123450",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty currency",
			request: MTNPaymentRequest{
				Amount:     "100",
				ExternalID: "test-external-id-123",
				Payer: MTNPayer{
					PartyIDType: "MSISDN",
					PartyID:     "46733123450",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty external ID",
			request: MTNPaymentRequest{
				Amount:   "100",
				Currency: "EUR",
				Payer: MTNPayer{
					PartyIDType: "MSISDN",
					PartyID:     "46733123450",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid payer party ID type",
			request: MTNPaymentRequest{
				Amount:     "100",
				Currency:   "EUR",
				ExternalID: "test-external-id-123",
				Payer: MTNPayer{
					PartyIDType: "INVALID",
					PartyID:     "46733123450",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.Amount)
				assert.NotEmpty(t, tt.request.Currency)
				assert.NotEmpty(t, tt.request.ExternalID)
				assert.NotEmpty(t, tt.request.Payer.PartyID)
				assert.Contains(t, []string{"MSISDN", "EMAIL", "PARTY_CODE"}, tt.request.Payer.PartyIDType)
			}
		})
	}
}

func TestMTNPayer_Validation(t *testing.T) {
	tests := []struct {
		name    string
		payer   MTNPayer
		wantErr bool
	}{
		{
			name: "Valid MSISDN payer",
			payer: MTNPayer{
				PartyIDType: "MSISDN",
				PartyID:     "46733123450",
			},
			wantErr: false,
		},
		{
			name: "Valid EMAIL payer",
			payer: MTNPayer{
				PartyIDType: "EMAIL",
				PartyID:     "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid PARTY_CODE payer",
			payer: MTNPayer{
				PartyIDType: "PARTY_CODE",
				PartyID:     "PARTY123",
			},
			wantErr: false,
		},
		{
			name: "Invalid party ID type",
			payer: MTNPayer{
				PartyIDType: "INVALID",
				PartyID:     "46733123450",
			},
			wantErr: true,
		},
		{
			name: "Empty party ID",
			payer: MTNPayer{
				PartyIDType: "MSISDN",
				PartyID:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NotEmpty(t, tt.payer.PartyID)
				assert.Contains(t, []string{"MSISDN", "EMAIL", "PARTY_CODE"}, tt.payer.PartyIDType)
			}
		})
	}
}

func TestOrangePaymentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request OrangePaymentRequest
		wantErr bool
	}{
		{
			name: "Valid Orange payment request",
			request: OrangePaymentRequest{
				MerchantKey:    "test-merchant-key",
				Currency:       "XOF",
				OrderID:        "order-123",
				Amount:         1000.0,
				ReturnURL:      "https://example.com/return",
				CancelURL:      "https://example.com/cancel",
				NotifURL:       "https://example.com/notify",
				Lang:           "fr",
				Reference:      "ref-123",
				CustomerMSISDN: "22507123456",
			},
			wantErr: false,
		},
		{
			name: "Empty merchant key",
			request: OrangePaymentRequest{
				Currency:  "XOF",
				OrderID:   "order-123",
				Amount:    1000.0,
				ReturnURL: "https://example.com/return",
				CancelURL: "https://example.com/cancel",
				NotifURL:  "https://example.com/notify",
			},
			wantErr: true,
		},
		{
			name: "Zero amount",
			request: OrangePaymentRequest{
				MerchantKey: "test-merchant-key",
				Currency:    "XOF",
				OrderID:     "order-123",
				Amount:      0,
				ReturnURL:   "https://example.com/return",
				CancelURL:   "https://example.com/cancel",
				NotifURL:    "https://example.com/notify",
			},
			wantErr: true,
		},
		{
			name: "Negative amount",
			request: OrangePaymentRequest{
				MerchantKey: "test-merchant-key",
				Currency:    "XOF",
				OrderID:     "order-123",
				Amount:      -100.0,
				ReturnURL:   "https://example.com/return",
				CancelURL:   "https://example.com/cancel",
				NotifURL:    "https://example.com/notify",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NotEmpty(t, tt.request.MerchantKey)
				assert.NotEmpty(t, tt.request.Currency)
				assert.NotEmpty(t, tt.request.OrderID)
				assert.Greater(t, tt.request.Amount, 0.0)
				assert.NotEmpty(t, tt.request.ReturnURL)
				assert.NotEmpty(t, tt.request.CancelURL)
				assert.NotEmpty(t, tt.request.NotifURL)
			}
		})
	}
}

func TestMobileMoneyError_Error(t *testing.T) {
	tests := []struct {
		name     string
		mmError  MobileMoneyError
		expected string
	}{
		{
			name: "MTN API error",
			mmError: MobileMoneyError{
				Code:    "INVALID_MSISDN",
				Message: "Invalid MSISDN format",
				Details: "The provided MSISDN does not match the expected format",
			},
			expected: "Invalid MSISDN format",
		},
		{
			name: "Orange API error",
			mmError: MobileMoneyError{
				Code:    "INSUFFICIENT_FUNDS",
				Message: "Insufficient funds in wallet",
			},
			expected: "Insufficient funds in wallet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mmError
			assert.Equal(t, tt.expected, err.Error())
			assert.Equal(t, tt.mmError.Code, err.Code)
			assert.Equal(t, tt.mmError.Message, err.Message)
		})
	}
}

// TestPaymentStatus_Constants is now in common_test.go

func TestMTNPaymentResponse_Fields(t *testing.T) {
	response := MTNPaymentResponse{
		ReferenceID: "ref-123-456",
		Status:      "PENDING",
		Reason:      "Payment initiated successfully",
	}

	assert.Equal(t, "ref-123-456", response.ReferenceID)
	assert.Equal(t, "PENDING", response.Status)
	assert.Equal(t, "Payment initiated successfully", response.Reason)
}

func TestMTNPaymentStatusResponse_Fields(t *testing.T) {
	response := MTNPaymentStatusResponse{
		Amount:                 "100",
		Currency:               "EUR",
		FinancialTransactionID: "ft-123-456",
		ExternalID:             "ext-123",
		Payer: MTNPayer{
			PartyIDType: "MSISDN",
			PartyID:     "46733123450",
		},
		PayerMessage: "Payment for services",
		PayeeNote:    "Thank you",
		Status:       "SUCCEEDED",
		Metadata: map[string]string{
			"orderId": "order-123",
		},
	}

	assert.Equal(t, "100", response.Amount)
	assert.Equal(t, "EUR", response.Currency)
	assert.Equal(t, "ft-123-456", response.FinancialTransactionID)
	assert.Equal(t, "SUCCEEDED", response.Status)
	assert.Equal(t, "MSISDN", response.Payer.PartyIDType)
	assert.Equal(t, "order-123", response.Metadata["orderId"])
}

func TestOrangePaymentResponse_Fields(t *testing.T) {
	response := OrangePaymentResponse{
		PayToken:    "pay-token-123",
		PaymentURL:  "https://payment.orange.com/pay/123",
		NotifToken:  "notif-token-456",
		MerchantKey: "merchant-key-789",
	}

	assert.Equal(t, "pay-token-123", response.PayToken)
	assert.Equal(t, "https://payment.orange.com/pay/123", response.PaymentURL)
	assert.Equal(t, "notif-token-456", response.NotifToken)
	assert.Equal(t, "merchant-key-789", response.MerchantKey)
}

func TestOrangePaymentStatusResponse_Fields(t *testing.T) {
	response := OrangePaymentStatusResponse{
		Status:      "SUCCESS",
		TxnID:       "txn-123-456",
		Amount:      1000.0,
		Currency:    "XOF",
		OrderID:     "order-123",
		MerchantKey: "merchant-key-789",
		Reference:   "ref-123",
		Message:     "Payment completed successfully",
	}

	assert.Equal(t, "SUCCESS", response.Status)
	assert.Equal(t, "txn-123-456", response.TxnID)
	assert.Equal(t, 1000.0, response.Amount)
	assert.Equal(t, "XOF", response.Currency)
	assert.Equal(t, "order-123", response.OrderID)
	assert.Equal(t, "Payment completed successfully", response.Message)
}
