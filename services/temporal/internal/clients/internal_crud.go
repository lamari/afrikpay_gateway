package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/models"
)

// CrudClient implements the models.CrudClient interface for interacting with the CRUD service
type CrudClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// CrudClientConfig holds the configuration for the CRUD client
type CrudClientConfig struct {
	BaseURL string
}

// NewCrudClient creates a new CRUD client
func NewCrudClient(config *config.CrudConfig) (*CrudClient, error) {
	if config.BaseURL == "" {
		return nil, models.NewClientError("INVALID_CONFIG", "base URL is required", false)
	}

	client := &CrudClient{
		baseURL:    config.BaseURL,
		httpClient: &http.Client{},
	}

	return client, nil
}

// HealthCheck checks if the CRUD service is healthy
func (c *CrudClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/health", c.baseURL), nil)
	if err != nil {
		return models.NewClientError("REQUEST_ERROR", err.Error(), true)
	}

	c.setAuthHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.NewClientError("CONNECTION_ERROR", fmt.Sprintf("failed to connect to CRUD service: %s", err.Error()), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.NewClientError("SERVICE_UNAVAILABLE", fmt.Sprintf("CRUD service returned status code %d", resp.StatusCode), true)
	}

	return nil
}

// GetName returns the name of the client
func (c *CrudClient) GetName() string {
	return "internal-crud"
}

// Close closes the client
func (c *CrudClient) Close() error {
	// Nothing to close for HTTP client
	return nil
}

// CreateTransaction creates a new transaction in the CRUD service
func (c *CrudClient) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.TransactionResponse, error) {
	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return nil, models.NewClientError("SERIALIZATION_ERROR", err.Error(), false)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/transactions", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, models.NewClientError("REQUEST_ERROR", err.Error(), false)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, models.NewClientError("CONNECTION_ERROR", fmt.Sprintf("failed to connect to CRUD service: %s", err.Error()), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error != "" {
			return nil, models.NewClientError("API_ERROR", errorResp.Error, false)
		}
		return nil, models.NewClientError("API_ERROR", fmt.Sprintf("unexpected status code: %d", resp.StatusCode), false)
	}

	var response models.TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, models.NewClientError("DESERIALIZATION_ERROR", err.Error(), false)
	}

	return &response, nil
}

// UpdateWalletBalance updates a wallet's balance in the CRUD service
func (c *CrudClient) UpdateWalletBalance(ctx context.Context, walletID string, amount float64, currency string) (*models.WalletResponse, error) {
	if walletID == "" {
		return nil, models.NewClientError("INVALID_WALLET", "wallet ID is required", false)
	}
	if amount == 0 {
		return nil, models.NewClientError("INVALID_AMOUNT", "amount cannot be zero", false)
	}
	if currency == "" {
		return nil, models.NewClientError("INVALID_CURRENCY", "currency is required", false)
	}

	updateRequest := struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}{
		Amount:   amount,
		Currency: currency,
	}

	jsonData, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, models.NewClientError("SERIALIZATION_ERROR", err.Error(), false)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("%s/api/v1/wallets/%s/balance", c.baseURL, walletID), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, models.NewClientError("REQUEST_ERROR", err.Error(), false)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, models.NewClientError("CONNECTION_ERROR", fmt.Sprintf("failed to connect to CRUD service: %s", err.Error()), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error != "" {
			return nil, models.NewClientError("API_ERROR", errorResp.Error, false)
		}
		return nil, models.NewClientError("API_ERROR", fmt.Sprintf("unexpected status code: %d", resp.StatusCode), false)
	}

	var response models.WalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, models.NewClientError("DESERIALIZATION_ERROR", err.Error(), false)
	}

	return &response, nil
}

// GetWallet retrieves wallet information by user ID and currency
func (c *CrudClient) GetWallet(ctx context.Context, userID string, currency string) (*models.WalletResponse, error) {
	if userID == "" {
		return nil, models.NewClientError("INVALID_USER", "user ID is required", false)
	}
	if currency == "" {
		return nil, models.NewClientError("INVALID_CURRENCY", "currency is required", false)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/users/%s/wallets?currency=%s", c.baseURL, userID, currency), nil)
	if err != nil {
		return nil, models.NewClientError("REQUEST_ERROR", err.Error(), false)
	}

	c.setAuthHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, models.NewClientError("CONNECTION_ERROR", fmt.Sprintf("failed to connect to CRUD service: %s", err.Error()), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, models.NewClientError("WALLET_NOT_FOUND", fmt.Sprintf("wallet not found for user %s with currency %s", userID, currency), false)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error != "" {
			return nil, models.NewClientError("API_ERROR", errorResp.Error, false)
		}
		return nil, models.NewClientError("API_ERROR", fmt.Sprintf("unexpected status code: %d", resp.StatusCode), false)
	}

	var response models.WalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, models.NewClientError("DESERIALIZATION_ERROR", err.Error(), false)
	}

	return &response, nil
}

// setAuthHeader sets the authorization header for requests
func (c *CrudClient) setAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
}
