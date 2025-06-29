package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/afrikpay/gateway/internal/models"
)

func TestOrangeClient_InitiatePayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentRequest *models.PaymentRequest
		mockResponse   string
		mockStatusCode int
		expectError    bool
		expectedStatus string
	}{
		{
			name: "Valid payment initiation",
			paymentRequest: &models.PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				ExternalID:  "ext-123",
				PhoneNumber: "22505123456",
				Description: "Payment for services",
				CallbackURL: "https://webhook.site/callback",
				Metadata: map[string]string{
					"orderId": "order-123",
				},
			},
			mockResponse: `{
				"transactionId": "orange-123-456",
				"status": "PENDING",
				"message": "Payment initiated successfully"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "PENDING",
		},
		{
			name: "Invalid phone number",
			paymentRequest: &models.PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				ExternalID:  "ext-123",
				PhoneNumber: "invalid-phone",
				Description: "Payment for services",
			},
			mockResponse:   `{"error":"INVALID_PHONE_NUMBER","message":"Invalid phone number format"}`,
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Insufficient balance",
			paymentRequest: &models.PaymentRequest{
				Amount:      10000.0,
				Currency:    "XOF",
				ExternalID:  "ext-124",
				PhoneNumber: "22505123456",
				Description: "Large payment",
			},
			mockResponse:   `{"error":"INSUFFICIENT_BALANCE","message":"Insufficient balance in customer account"}`,
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Service unavailable",
			paymentRequest: &models.PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				ExternalID:  "ext-125",
				PhoneNumber: "22505123456",
				Description: "Payment during maintenance",
			},
			mockResponse:   `{"error":"SERVICE_UNAVAILABLE","message":"Service temporarily unavailable"}`,
			mockStatusCode: http.StatusServiceUnavailable,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/omcoreapis/1.0.2/mp/pay", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.OrangeConfig{
				BaseURL:       server.URL,
				ClientID:      "test-api-key",
				ClientSecret:  "test-api-secret",
				Authorization: "test-auth",
				Timeout:       30 * time.Second,
			}

			client := NewOrangeClient(config)
			ctx := context.Background()

			// Execute test
			response, err := client.InitiatePayment(ctx, tt.paymentRequest)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.NotEmpty(t, response.ReferenceID)
				assert.Equal(t, tt.expectedStatus, string(response.Status))
				assert.NotEmpty(t, response.Message)
				assert.NotEmpty(t, response.PaymentURL)
			}
		})
	}
}

func TestOrangeClient_GetPaymentStatus(t *testing.T) {
	tests := []struct {
		name           string
		referenceID    string
		mockResponse   string
		mockStatusCode int
		expectError    bool
		expectedStatus string
	}{
		{
			name:        "Orange payment successful",
			referenceID: "orange-123-456",
			mockResponse: `{
				"transactionId": "orange-123-456",
				"externalId": "ext-123",
				"status": "SUCCESS",
				"amount": 100.0,
				"currency": "XOF",
				"phoneNumber": "22505123456",
				"message": "Payment completed successfully",
				"timestamp": "2023-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "SUCCESS",
		},
		{
			name:        "Orange payment pending",
			referenceID: "orange-123-457",
			mockResponse: `{
				"transactionId": "orange-123-457",
				"externalId": "ext-124",
				"status": "PENDING",
				"amount": 200.0,
				"currency": "XOF",
				"phoneNumber": "22505123457",
				"message": "Payment is being processed",
				"timestamp": "2023-01-01T12:01:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "PENDING",
		},
		{
			name:        "Orange payment failed",
			referenceID: "orange-123-458",
			mockResponse: `{
				"transactionId": "orange-123-458",
				"externalId": "ext-125",
				"status": "FAILED",
				"amount": 300.0,
				"currency": "XOF",
				"phoneNumber": "22505123458",
				"message": "Payment rejected by user",
				"timestamp": "2023-01-01T12:02:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "FAILED",
		},
		{
			name:           "Orange payment not found",
			referenceID:    "non-existent",
			mockResponse:   `{"error":"NOT_FOUND","message":"Transaction not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/omcoreapis/1.0.2/mp/pay/"+tt.referenceID, r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))

				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.OrangeConfig{
				BaseURL:       server.URL,
				ClientID:      "test-api-key",
				ClientSecret:  "test-api-secret",
				Authorization: "test-auth",
				Timeout:       30 * time.Second,
			}

			client := NewOrangeClient(config)
			ctx := context.Background()

			// Execute test
			status, err := client.GetPaymentStatus(ctx, tt.referenceID)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, status)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, status)
				assert.Equal(t, tt.referenceID, status.ReferenceID)
				assert.Equal(t, tt.expectedStatus, string(status.Status))
			}
		})
	}
}

// Mock implementation for Orange testing
type MockOrangeClient struct {
	payments    map[string]*models.PaymentStatusResponse
	healthError error
}

func NewMockOrangeClient() *MockOrangeClient {
	return &MockOrangeClient{
		payments: make(map[string]*models.PaymentStatusResponse),
	}
}

func (m *MockOrangeClient) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	referenceID := "mock-orange-ref-" + request.ExternalID

	// Store payment status for later retrieval
	m.payments[referenceID] = &models.PaymentStatusResponse{
		ReferenceID: referenceID,
		ExternalID:  request.ExternalID,
		Status:      models.PaymentStatusPending,
		Amount:      request.Amount,
		Currency:    request.Currency,
		PhoneNumber: request.PhoneNumber,
		Message:     "Payment initiated successfully",
		Timestamp:   time.Now(),
	}

	return &models.PaymentResponse{
		ReferenceID: referenceID,
		Status:      models.PaymentStatusPending,
		PaymentURL:  "https://mock-orange.com/pay/" + referenceID,
		Message:     "Payment initiated successfully",
	}, nil
}

func (m *MockOrangeClient) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	payment, exists := m.payments[referenceID]
	if !exists {
		return nil, models.ClientError{
			Code:    "PAYMENT_NOT_FOUND",
			Message: "Payment not found",
		}
	}

	return payment, nil
}

func (m *MockOrangeClient) HealthCheck(ctx context.Context) error {
	return m.healthError
}

func (m *MockOrangeClient) GetName() string {
	return "mock-orange"
}

func (m *MockOrangeClient) SetPaymentStatus(referenceID string, status models.PaymentStatus) {
	if payment, exists := m.payments[referenceID]; exists {
		payment.Status = status
		if status == models.PaymentStatusSuccess {
			payment.TransactionID = "mock-txn-" + referenceID
			payment.Message = "Payment completed successfully"
		} else if status == models.PaymentStatusFailed {
			payment.Message = "Payment failed"
		}
	}
}

func (m *MockOrangeClient) SetHealthError(err error) {
	m.healthError = err
}

func TestMockOrangeClient(t *testing.T) {
	mock := NewMockOrangeClient()
	ctx := context.Background()

	// Test InitiatePayment
	request := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "XOF",
		ExternalID:  "ext-mock-123",
		PhoneNumber: "22505123456",
		Description: "Mock payment test",
	}

	response, err := mock.InitiatePayment(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, models.PaymentStatusPending, response.Status)

	// Test GetPaymentStatus
	status, err := mock.GetPaymentStatus(ctx, response.ReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, models.PaymentStatusPending, status.Status)

	// Test SetPaymentStatus
	mock.SetPaymentStatus(response.ReferenceID, models.PaymentStatusSuccess)
	status, err = mock.GetPaymentStatus(ctx, response.ReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, models.PaymentStatusSuccess, status.Status)

	// Test HealthCheck
	err = mock.HealthCheck(ctx)
	assert.NoError(t, err)

	// Test GetName
	assert.Equal(t, "mock-orange", mock.GetName())
}
