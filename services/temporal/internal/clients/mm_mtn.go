package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/afrikpay/gateway/internal/models"
)

// MTNClient implements the MobileMoneyClient interface for MTN Mobile Money
type MTNClient struct {
	config     models.MTNConfig
	httpClient *http.Client
}

// NewMTNClient creates a new MTN Mobile Money client
func NewMTNClient(config models.MTNConfig) *MTNClient {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &MTNClient{
		config:     config,
		httpClient: httpClient,
	}
}

// MTN Client Methods

// InitiatePayment initiates a payment with MTN Mobile Money
func (c *MTNClient) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	return c.initiatePayment(ctx, request)
}

// GetPaymentStatus gets the status of a payment
func (c *MTNClient) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	return c.getPaymentStatus(ctx, referenceID)
}

// HealthCheck performs a health check
func (c *MTNClient) HealthCheck(ctx context.Context) error {
	return c.healthCheck(ctx)
}

// GetName returns the client name
func (c *MTNClient) GetName() string {
	return "mtn"
}

// Close closes the client
func (c *MTNClient) Close() error {
	return nil
}

// initiatePayment initiates payment with MTN API
func (c *MTNClient) initiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", c.config.BaseURL)

	payload := map[string]interface{}{
		"amount":     fmt.Sprintf("%.2f", request.Amount),
		"currency":   request.Currency,
		"externalId": request.ExternalID,
		"payer": map[string]string{
			"partyIdType": "MSISDN",
			"partyId":     request.PhoneNumber,
		},
		"payerMessage": request.Description,
		"payeeNote":    request.Description,
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
	req.Header.Set("Authorization", "Bearer "+c.config.PrimaryKey)
	req.Header.Set("X-Reference-Id", fmt.Sprintf("mtn-%d", time.Now().UnixNano()))
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("MTN API error: status %d", resp.StatusCode)
	}

	referenceID := resp.Header.Get("X-Reference-Id")
	if referenceID == "" {
		referenceID = fmt.Sprintf("mtn-%d", time.Now().UnixNano())
	}

	return &models.PaymentResponse{
		ReferenceID: referenceID,
		Status:      models.PaymentStatusPending,
		PaymentURL:  fmt.Sprintf("%s/pay/%s", c.config.BaseURL, referenceID),
		Message:     "Payment initiated successfully",
	}, nil
}

// getPaymentStatus gets payment status from MTN API
func (c *MTNClient) getPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", c.config.BaseURL, referenceID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.config.PrimaryKey)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MTN API error: status %d", resp.StatusCode)
	}

	var mtnResp struct {
		Amount                 string `json:"amount"`
		Currency               string `json:"currency"`
		ExternalID             string `json:"externalId"`
		Status                 string `json:"status"`
		Reason                 string `json:"reason"`
		FinancialTransactionID string `json:"financialTransactionId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mtnResp); err != nil {
		return nil, err
	}

	// Convert MTN status to our status
	var status models.PaymentStatus
	switch mtnResp.Status {
	case "SUCCESSFUL":
		status = models.PaymentStatusSuccess
	case "FAILED":
		status = models.PaymentStatusFailed
	case "PENDING":
		status = models.PaymentStatusPending
	default:
		status = models.PaymentStatusPending
	}

	return &models.PaymentStatusResponse{
		ReferenceID: referenceID,
		ExternalID:  mtnResp.ExternalID,
		Status:      status,
		Amount: func() float64 {
			if f, err := strconv.ParseFloat(mtnResp.Amount, 64); err == nil {
				return f
			}
			return 0.0
		}(),
		Currency:      mtnResp.Currency,
		TransactionID: mtnResp.FinancialTransactionID,
		Message:       mtnResp.Reason,
		Timestamp:     time.Now(),
	}, nil
}

// healthCheck performs health check against MTN API
func (c *MTNClient) healthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/collection/v1_0/accountbalance", c.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.config.PrimaryKey)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MTN health check failed: status %d", resp.StatusCode)
	}

	return nil
}
