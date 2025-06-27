package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/afrikpay/gateway/services/client/internal/models"
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
				"message": "Payment initiated successfully"
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
			config := models.MobileMoneyConfig{
				BaseURL:         server.URL,
				APIKey:          "test-api-key",
				APISecret:       "test-api-secret",
				SubscriptionKey: "test-subscription-key",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
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
			name:        "Payment successful",
			referenceID: "ref-mtn-123-456",
			mockResponse: `{
				"referenceId": "ref-mtn-123-456",
				"externalId": "ext-123",
				"status": "SUCCEEDED",
				"amount": 100.0,
				"currency": "XOF",
				"phoneNumber": "22507123456",
				"transactionId": "mtn-txn-789",
				"message": "Payment completed successfully",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "SUCCEEDED",
		},
		{
			name:        "Payment failed",
			referenceID: "ref-mtn-123-457",
			mockResponse: `{
				"referenceId": "ref-mtn-123-457",
				"externalId": "ext-124",
				"status": "FAILED",
				"amount": 100.0,
				"currency": "XOF",
				"phoneNumber": "22507123456",
				"message": "Payment failed due to insufficient balance",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "FAILED",
		},
		{
			name:           "Payment not found",
			referenceID:    "ref-not-found",
			mockResponse:   `{"error":"PAYMENT_NOT_FOUND","message":"Payment reference not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:        "Payment pending",
			referenceID: "ref-mtn-123-458",
			mockResponse: `{
				"referenceId": "ref-mtn-123-458",
				"externalId": "ext-125",
				"status": "PENDING",
				"amount": 100.0,
				"currency": "XOF",
				"phoneNumber": "22507123456",
				"message": "Payment is being processed",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "PENDING",
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
				assert.NotEmpty(t, r.Header.Get("X-Target-Environment"))
				
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MobileMoneyConfig{
				BaseURL:         server.URL,
				APIKey:          "test-api-key",
				APISecret:       "test-api-secret",
				SubscriptionKey: "test-subscription-key",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
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
			name: "Valid Orange payment initiation",
			paymentRequest: &models.PaymentRequest{
				Amount:      50.0,
				Currency:    "XOF",
				ExternalID:  "ext-orange-123",
				PhoneNumber: "22507654321",
				Description: "Orange Money payment",
				CallbackURL: "https://webhook.site/orange-callback",
			},
			mockResponse: `{
				"referenceId": "ref-orange-123-456",
				"status": "PENDING",
				"paymentUrl": "https://orange.com/pay/123",
				"message": "Payment initiated successfully"
			}`,
			mockStatusCode: http.StatusAccepted,
			expectError:    false,
			expectedStatus: "PENDING",
		},
		{
			name: "Invalid amount",
			paymentRequest: &models.PaymentRequest{
				Amount:      -10.0,
				Currency:    "XOF",
				ExternalID:  "ext-orange-124",
				PhoneNumber: "22507654321",
				Description: "Invalid amount payment",
			},
			mockResponse:   `{"error":"INVALID_AMOUNT","message":"Amount must be positive"}`,
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Customer not found",
			paymentRequest: &models.PaymentRequest{
				Amount:      100.0,
				Currency:    "XOF",
				ExternalID:  "ext-orange-125",
				PhoneNumber: "22500000000",
				Description: "Payment to non-existent customer",
			},
			mockResponse:   `{"error":"CUSTOMER_NOT_FOUND","message":"Customer not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/orange-money-webpay/dev/v1/webpayment", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MobileMoneyConfig{
				BaseURL:         server.URL,
				APIKey:          "test-orange-key",
				APISecret:       "test-orange-secret",
				SubscriptionKey: "test-orange-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
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
			referenceID: "ref-orange-123-456",
			mockResponse: `{
				"referenceId": "ref-orange-123-456",
				"externalId": "ext-orange-123",
				"status": "SUCCEEDED",
				"amount": 50.0,
				"currency": "XOF",
				"phoneNumber": "22507654321",
				"transactionId": "orange-txn-789",
				"message": "Payment completed successfully",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "SUCCEEDED",
		},
		{
			name:        "Orange payment timeout",
			referenceID: "ref-orange-123-457",
			mockResponse: `{
				"referenceId": "ref-orange-123-457",
				"externalId": "ext-orange-124",
				"status": "TIMEOUT",
				"amount": 50.0,
				"currency": "XOF",
				"phoneNumber": "22507654321",
				"message": "Payment timed out",
				"timestamp": "2024-01-01T12:00:00Z"
			}`,
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedStatus: "TIMEOUT",
		},
		{
			name:           "Orange payment not found",
			referenceID:    "ref-orange-not-found",
			mockResponse:   `{"error":"PAYMENT_NOT_FOUND","message":"Payment reference not found"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/orange-money-webpay/dev/v1/webpayment/" + tt.referenceID
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				assert.NotEmpty(t, r.Header.Get("Authorization"))
				
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MobileMoneyConfig{
				BaseURL:         server.URL,
				APIKey:          "test-orange-key",
				APISecret:       "test-orange-secret",
				SubscriptionKey: "test-orange-subscription",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			}

			client := NewOrangeClient(config)
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

func TestMobileMoneyClient_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		clientType     string
		mockStatusCode int
		expectError    bool
	}{
		{
			name:           "MTN healthy service",
			clientType:     "mtn",
			mockStatusCode: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Orange healthy service",
			clientType:     "orange",
			mockStatusCode: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "MTN service unavailable",
			clientType:     "mtn",
			mockStatusCode: http.StatusServiceUnavailable,
			expectError:    true,
		},
		{
			name:           "Orange service unavailable",
			clientType:     "orange",
			mockStatusCode: http.StatusServiceUnavailable,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var expectedPath string
				if tt.clientType == "mtn" {
					expectedPath = "/collection/v1_0/account/balance"
				} else {
					expectedPath = "/orange-money-webpay/dev/v1/account/balance"
				}
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockStatusCode == http.StatusOK {
					w.Write([]byte(`{"balance": "1000.00", "currency": "XOF"}`))
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			config := models.MobileMoneyConfig{
				BaseURL:         server.URL,
				APIKey:          "test-api-key",
				APISecret:       "test-api-secret",
				SubscriptionKey: "test-subscription-key",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
			}

			var client models.MobileMoneyClient
			if tt.clientType == "mtn" {
				client = NewMTNClient(config)
			} else {
				client = NewOrangeClient(config)
			}

			ctx := context.Background()

			// Execute test
			err := client.HealthCheck(ctx)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMobileMoneyClient_GetName(t *testing.T) {
	config := models.MobileMoneyConfig{
		BaseURL:         "https://api.example.com",
		APIKey:          "test-api-key",
		APISecret:       "test-api-secret",
		SubscriptionKey: "test-subscription-key",
		Timeout:         30 * time.Second,
		MaxRetries:      3,
	}

	mtnClient := NewMTNClient(config)
	assert.Equal(t, "mtn", mtnClient.GetName())

	orangeClient := NewOrangeClient(config)
	assert.Equal(t, "orange", orangeClient.GetName())
}

func TestMobileMoneyClient_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"referenceId":"ref-123","status":"PENDING"}`))
	}))
	defer server.Close()

	// Create client with short timeout
	config := models.MobileMoneyConfig{
		BaseURL:         server.URL,
		APIKey:          "test-api-key",
		APISecret:       "test-api-secret",
		SubscriptionKey: "test-subscription-key",
		Timeout:         1 * time.Second,
		MaxRetries:      1,
	}

	client := NewMTNClient(config)
	ctx := context.Background()

	paymentRequest := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "XOF",
		ExternalID:  "ext-timeout-test",
		PhoneNumber: "22507123456",
		Description: "Timeout test payment",
	}

	// Execute test - should timeout
	_, err := client.InitiatePayment(ctx, paymentRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestMobileMoneyClient_RetryLogic(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			// Fail first 2 calls
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"INTERNAL_ERROR","message":"Internal server error"}`))
		} else {
			// Succeed on 3rd call
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"referenceId":"ref-retry-123","status":"PENDING","message":"Payment initiated"}`))
		}
	}))
	defer server.Close()

	config := models.MobileMoneyConfig{
		BaseURL:         server.URL,
		APIKey:          "test-api-key",
		APISecret:       "test-api-secret",
		SubscriptionKey: "test-subscription-key",
		Timeout:         30 * time.Second,
		MaxRetries:      3,
	}

	client := NewMTNClient(config)
	ctx := context.Background()

	paymentRequest := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "XOF",
		ExternalID:  "ext-retry-test",
		PhoneNumber: "22507123456",
		Description: "Retry test payment",
	}

	// Execute test - should succeed after retries
	response, err := client.InitiatePayment(ctx, paymentRequest)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, callCount) // Should have made 3 calls
}

// Mock implementations for testing
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
