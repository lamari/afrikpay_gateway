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

func TestMTNClient_InitiatePayment(t *testing.T) {
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
				PhoneNumber: "22507123456",
				Description: "Payment for services",
				CallbackURL: "https://webhook.site/callback",
				Metadata: map[string]string{
					"orderId": "order-123",
				},
			},
			mockResponse: `{
				"referenceId": "ref-mtn-123-456",
				"status": "PENDING",
				"paymentUrl": "https://mtn.com/pay/123",
				"reason": "Payment initiated successfully"
			}`,
			mockStatusCode: http.StatusAccepted,
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
				PhoneNumber: "22507123456",
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
				PhoneNumber: "22507123456",
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
				assert.Equal(t, "/collection/v1_0/requesttopay", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))
				assert.NotEmpty(t, r.Header.Get("X-Reference-Id"))
				assert.NotEmpty(t, r.Header.Get("X-Target-Environment"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MTNConfig{
				BaseURL:      server.URL,
				PrimaryKey:   "test-api-key",
				SecondaryKey: "test-api-secret",
				Timeout:      30 * time.Second,
				MaxRetries:   3,
			}

			client := NewMTNClient(config)
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
			}
		})
	}
}

func TestMTNClient_GetPaymentStatus(t *testing.T) {
	tests := []struct {
		name           string
		referenceID    string
		mockResponse   string
		mockStatusCode int
		expectError    bool
		expectedStatus string
	}{
		{
			name:        "MTN payment successful",
			referenceID: "ref-mtn-123-456",
			mockResponse: `{
				"referenceId": "ref-mtn-123-456",
				"externalId": "ext-mtn-123",
				"status": "SUCCESS",
				"amount": 50.0,
				"currency": "XOF",
				"phoneNumber": "22507123456",
				"transactionId": "mtn-txn-789",
				"message": "Payment completed successfully",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "SUCCESS",
		},
		{
			name:        "MTN payment timeout",
			referenceID: "ref-mtn-123-457",
			mockResponse: `{
				"referenceId": "ref-mtn-123-457",
				"externalId": "ext-mtn-124",
				"status": "PENDING",
				"amount": 50.0,
				"currency": "XOF",
				"phoneNumber": "22507123456",
				"message": "Payment timed out",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "PENDING",
		},
		{
			name:           "MTN payment not found",
			referenceID:    "ref-mtn-not-found",
			mockResponse:   `{"error":"PAYMENT_NOT_FOUND","message":"Payment reference not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/collection/v1_0/requesttopay/" + tt.referenceID
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))

				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MTNConfig{
				BaseURL:      server.URL,
				PrimaryKey:   "test-mtn-key",
				SecondaryKey: "test-mtn-secret",
				Timeout:      30 * time.Second,
				MaxRetries:   3,
			}

			client := NewMTNClient(config)
			ctx := context.Background()

			// Execute test
			response, err := client.GetPaymentStatus(ctx, tt.referenceID)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, tt.referenceID, response.ReferenceID)
				assert.Equal(t, tt.expectedStatus, string(response.Status))
				assert.NotEmpty(t, response.Message)
			}
		})
	}
}

// Mock implementation for MTN testing
type MockMTNClient struct {
	payments    map[string]*models.PaymentStatusResponse
	healthError error
}

func NewMockMTNClient() *MockMTNClient {
	return &MockMTNClient{
		payments: make(map[string]*models.PaymentStatusResponse),
	}
}

func (m *MockMTNClient) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	referenceID := "mock-mtn-ref-" + request.ExternalID

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
		PaymentURL:  "https://mock-mtn.com/pay/" + referenceID,
		Message:     "Payment initiated successfully",
	}, nil
}

func (m *MockMTNClient) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	payment, exists := m.payments[referenceID]
	if !exists {
		return nil, models.ClientError{
			Code:    "PAYMENT_NOT_FOUND",
			Message: "Payment not found",
		}
	}

	return payment, nil
}

func (m *MockMTNClient) HealthCheck(ctx context.Context) error {
	return m.healthError
}

func (m *MockMTNClient) GetName() string {
	return "mock-mtn"
}

func (m *MockMTNClient) SetPaymentStatus(referenceID string, status models.PaymentStatus) {
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

func (m *MockMTNClient) SetHealthError(err error) {
	m.healthError = err
}

func TestMockMTNClient(t *testing.T) {
	mock := NewMockMTNClient()
	ctx := context.Background()

	// Test InitiatePayment
	request := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "XOF",
		ExternalID:  "ext-mock-123",
		PhoneNumber: "22507123456",
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
	assert.Equal(t, "mock-mtn", mock.GetName())
}
