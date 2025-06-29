package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCrudClient(t *testing.T) {
	// Test cases following the Given-When-Then structure
	tests := []struct {
		name          string
		config        *CrudClientConfig
		expectedError bool
	}{
		{
			name: "Valid configuration",
			config: &CrudClientConfig{
				BaseURL: "http://localhost:8002",
			},
			expectedError: false,
		},
		{
			name: "Missing base URL",
			config: &CrudClientConfig{
				BaseURL: "",
			},
			expectedError: true,
		},
		{
			name: "Missing auth token",
			config: &CrudClientConfig{
				BaseURL: "http://localhost:8002",
			},
			expectedError: true,
		},
		{
			name: "Default timeout applied",
			config: &CrudClientConfig{
				BaseURL: "http://localhost:8002",
			},
			expectedError: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When
			client, err := NewCrudClient(&config.CrudConfig{
				BaseURL: tc.config.BaseURL,
			})

			// Then
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tc.config.BaseURL, client.baseURL)
			}
		})
	}
}

func TestCrudClientHealthCheck(t *testing.T) {
	// Given
	tests := []struct {
		name           string
		statusCode     int
		expectedError  bool
		serverResponse string
	}{
		{
			name:           "Healthy service",
			statusCode:     http.StatusOK,
			expectedError:  false,
			serverResponse: `{"status":"ok"}`,
		},
		{
			name:           "Unhealthy service",
			statusCode:     http.StatusInternalServerError,
			expectedError:  true,
			serverResponse: `{"status":"error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				assert.Equal(t, "/health", r.URL.Path)
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				// Send response
				w.WriteHeader(tc.statusCode)
				fmt.Fprintln(w, tc.serverResponse)
			}))
			defer server.Close()

			// Create client with test server URL
			client, err := NewCrudClient(&config.CrudConfig{
				BaseURL: server.URL,
			})
			require.NoError(t, err)
			require.NotNil(t, client)

			// When
			err = client.HealthCheck(context.Background())

			// Then
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCrudClientCreateTransaction(t *testing.T) {
	// Given
	transaction := &models.Transaction{
		Type:     models.TransactionTypeDeposit,
		Status:   models.TransactionStatusPending,
		UserID:   "user123",
		Amount:   100.0,
		Currency: "USDT",
	}

	expectedResponse := &models.TransactionResponse{
		TransactionID: "tx123",
		UserID:        "user123",
		WalletID:      "wallet123",
		Type:          models.TransactionTypeDeposit,
		Status:        models.TransactionStatusPending,
		Amount:        100.0,
		Currency:      "USDT",
		CreatedAt:     time.Now(),
		Success:       true,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/api/v1/transactions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Check request body
		var reqTransaction models.Transaction
		err := json.NewDecoder(r.Body).Decode(&reqTransaction)
		assert.NoError(t, err)
		assert.Equal(t, transaction.UserID, reqTransaction.UserID)
		assert.Equal(t, transaction.Amount, reqTransaction.Amount)
		assert.Equal(t, transaction.Currency, reqTransaction.Currency)

		// Send response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewCrudClient(&config.CrudConfig{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// When
	response, err := client.CreateTransaction(context.Background(), transaction)

	// Then
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, expectedResponse.TransactionID, response.TransactionID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.WalletID, response.WalletID)
	assert.Equal(t, expectedResponse.Amount, response.Amount)
	assert.Equal(t, expectedResponse.Currency, response.Currency)
	assert.Equal(t, expectedResponse.Success, response.Success)
}

func TestCrudClientUpdateWalletBalance(t *testing.T) {
	// Given
	walletID := "wallet123"
	amount := 50.0
	currency := "USDT"

	expectedResponse := &models.WalletResponse{
		WalletID:  walletID,
		UserID:    "user123",
		Balance:   150.0, // Updated balance
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Success:   true,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, fmt.Sprintf("/api/v1/wallets/%s/balance", walletID), r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Check request body
		var reqUpdate struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		}
		err := json.NewDecoder(r.Body).Decode(&reqUpdate)
		assert.NoError(t, err)
		assert.Equal(t, amount, reqUpdate.Amount)
		assert.Equal(t, currency, reqUpdate.Currency)

		// Send response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewCrudClient(&config.CrudConfig{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// When
	response, err := client.UpdateWalletBalance(context.Background(), walletID, amount, currency)

	// Then
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, expectedResponse.WalletID, response.WalletID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Balance, response.Balance)
	assert.Equal(t, expectedResponse.Currency, response.Currency)
	assert.Equal(t, expectedResponse.Success, response.Success)
}

func TestCrudClientGetWallet(t *testing.T) {
	// Given
	userID := "user123"
	currency := "USDT"

	expectedResponse := &models.WalletResponse{
		WalletID:  "wallet123",
		UserID:    userID,
		Balance:   100.0,
		Currency:  currency,
		CreatedAt: time.Now(),
		Success:   true,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, fmt.Sprintf("/api/v1/users/%s/wallets", userID), r.URL.Path)
		assert.Equal(t, currency, r.URL.Query().Get("currency"))
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewCrudClient(&config.CrudConfig{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// When
	response, err := client.GetWallet(context.Background(), userID, currency)

	// Then
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, expectedResponse.WalletID, response.WalletID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Balance, response.Balance)
	assert.Equal(t, expectedResponse.Currency, response.Currency)
	assert.Equal(t, expectedResponse.Success, response.Success)
}

func TestCrudClientGetWalletNotFound(t *testing.T) {
	// Given
	userID := "user123"
	currency := "USDT"

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, fmt.Sprintf("/api/v1/users/%s/wallets", userID), r.URL.Path)
		assert.Equal(t, currency, r.URL.Query().Get("currency"))
		assert.Equal(t, "GET", r.Method)

		// Send not found response
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Wallet not found",
		})
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewCrudClient(&config.CrudConfig{
		BaseURL: server.URL,
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// When
	response, err := client.GetWallet(context.Background(), userID, currency)

	// Then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "wallet not found")
}

func TestCrudClientValidations(t *testing.T) {
	// Given
	client, err := NewCrudClient(&config.CrudConfig{
		BaseURL: "http://localhost:8002",
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test CreateTransaction validations
	t.Run("CreateTransaction_InvalidTransaction", func(t *testing.T) {
		// Test with invalid transaction
		invalidTx := &models.Transaction{
			// Missing required fields
		}
		response, err := client.CreateTransaction(context.Background(), invalidTx)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	// Test UpdateWalletBalance validations
	t.Run("UpdateWalletBalance_InvalidParams", func(t *testing.T) {
		// Test with missing wallet ID
		response, err := client.UpdateWalletBalance(context.Background(), "", 50.0, "USDT")
		assert.Error(t, err)
		assert.Nil(t, response)

		// Test with zero amount
		response, err = client.UpdateWalletBalance(context.Background(), "wallet123", 0, "USDT")
		assert.Error(t, err)
		assert.Nil(t, response)

		// Test with missing currency
		response, err = client.UpdateWalletBalance(context.Background(), "wallet123", 50.0, "")
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	// Test GetWallet validations
	t.Run("GetWallet_InvalidParams", func(t *testing.T) {
		// Test with missing user ID
		response, err := client.GetWallet(context.Background(), "", "USDT")
		assert.Error(t, err)
		assert.Nil(t, response)

		// Test with missing currency
		response, err = client.GetWallet(context.Background(), "user123", "")
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}
