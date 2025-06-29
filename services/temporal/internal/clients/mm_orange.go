package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/afrikpay/gateway/internal/models"
)

// OrangeClient implements the MobileMoneyClient interface for Orange Money
type OrangeClient struct {
	config     models.OrangeConfig
	httpClient *http.Client
}

// NewOrangeClient creates a new Orange Money client
func NewOrangeClient(config models.OrangeConfig) *OrangeClient {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &OrangeClient{
		config:     config,
		httpClient: httpClient,
	}
}

// Orange Client Methods

// InitiatePayment initiates a payment with Orange Money
func (c *OrangeClient) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	return c.initiatePayment(ctx, request)
}

// GetPaymentStatus gets the status of a payment
func (c *OrangeClient) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	return c.getPaymentStatus(ctx, referenceID)
}

// HealthCheck performs a health check
func (c *OrangeClient) HealthCheck(ctx context.Context) error {
	return c.healthCheck(ctx)
}

// GetName returns the client name
func (c *OrangeClient) GetName() string {
	return "orange"
}

// Close closes the client
func (c *OrangeClient) Close() error {
	return nil
}

// initiatePayment initiates payment with Orange API
func (c *OrangeClient) initiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	url := fmt.Sprintf("%s/omcoreapis/1.0.2/mp/pay", c.config.BaseURL)

	payload := map[string]interface{}{
		"partner": map[string]string{
			"idType": "MSISDN",
			"id":     request.PhoneNumber,
		},
		"amount": map[string]interface{}{
			"value": request.Amount,
			"unit":  request.Currency,
		},
		"reference":   request.ExternalID,
		"description": request.Description,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.config.Authorization)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Orange API error: status %d", resp.StatusCode)
	}

	var orangeResp struct {
		TransactionID string `json:"transactionId"`
		Status        string `json:"status"`
		Message       string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orangeResp); err != nil {
		return nil, err
	}

	return &models.PaymentResponse{
		ReferenceID: orangeResp.TransactionID,
		Status:      models.PaymentStatusPending,
		PaymentURL:  fmt.Sprintf("%s/pay/%s", c.config.BaseURL, orangeResp.TransactionID),
		Message:     orangeResp.Message,
	}, nil
}

// getPaymentStatus gets payment status from Orange API
func (c *OrangeClient) getPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/omcoreapis/1.0.2/mp/pay/%s", c.config.BaseURL, referenceID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.config.Authorization)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Orange API error: status %d", resp.StatusCode)
	}

	var orangeResp struct {
		TransactionID string  `json:"transactionId"`
		Status        string  `json:"status"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		Reference     string  `json:"reference"`
		Message       string  `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orangeResp); err != nil {
		return nil, err
	}

	// Convert Orange status to our status
	var status models.PaymentStatus
	switch orangeResp.Status {
	case "SUCCESS", "COMPLETED":
		status = models.PaymentStatusSuccess
	case "FAILED", "ERROR":
		status = models.PaymentStatusFailed
	case "PENDING", "PROCESSING":
		status = models.PaymentStatusPending
	default:
		status = models.PaymentStatusPending
	}

	return &models.PaymentStatusResponse{
		ReferenceID:   referenceID,
		ExternalID:    orangeResp.Reference,
		Status:        status,
		Amount:        orangeResp.Amount,
		Currency:      orangeResp.Currency,
		TransactionID: orangeResp.TransactionID,
		Message:       orangeResp.Message,
		Timestamp:     time.Now(),
	}, nil
}

// healthCheck performs health check against Orange API
func (c *OrangeClient) healthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/omcoreapis/1.0.2/mp/balance", c.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.config.Authorization)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Orange health check failed: status %d", resp.StatusCode)
	}

	return nil
}
