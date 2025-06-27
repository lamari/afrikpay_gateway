package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/models"
)

// MTNClient implements the MobileMoneyClient interface for MTN Mobile Money
type MTNClient struct {
	config     models.MobileMoneyConfig
	httpClient *http.Client
}

// OrangeClient implements the MobileMoneyClient interface for Orange Money
type OrangeClient struct {
	config     models.MobileMoneyConfig
	httpClient *http.Client
}

// NewMobileMoneyClient creates a new mobile money client based on provider
func NewMobileMoneyClient(config models.MobileMoneyConfig) (models.MobileMoneyClient, error) {
	switch config.Provider {
	case "mtn":
		return NewMTNClient(config), nil
	case "orange":
		return NewOrangeClient(config), nil
	default:
		return nil, fmt.Errorf("unsupported mobile money provider: %s", config.Provider)
	}
}

// NewMTNClient creates a new MTN Mobile Money client
func NewMTNClient(config models.MobileMoneyConfig) *MTNClient {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &MTNClient{
		config:     config,
		httpClient: httpClient,
	}
}

// NewOrangeClient creates a new Orange Money client
func NewOrangeClient(config models.MobileMoneyConfig) *OrangeClient {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &OrangeClient{
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

// GetResilienceStats returns resilience statistics
// Deprecated: Resilience is now handled by Temporal workflows
func (c *MTNClient) GetResilienceStats() *models.ResilienceStats {
	// Return empty stats as resilience is handled by Temporal
	return &models.ResilienceStats{
		TotalRequests:      0,
		SuccessfulRequests: 0,
		FailedRequests:     0,
		LastReset:          time.Now(),
	}
}

// ResetResilienceStats resets resilience statistics
// Deprecated: Resilience is now handled by Temporal workflows
func (c *MTNClient) ResetResilienceStats() {
	// No-op as resilience is handled by Temporal
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
		"payerMessage":  request.Description,
		"payeeNote":     request.Description,
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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("X-Reference-Id", fmt.Sprintf("mtn-%d", time.Now().UnixNano()))
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SubscriptionKey)

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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SubscriptionKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MTN API error: status %d", resp.StatusCode)
	}

	var mtnResp struct {
		Amount         string `json:"amount"`
		Currency       string `json:"currency"`
		ExternalID     string `json:"externalId"`
		Status         string `json:"status"`
		Reason         string `json:"reason"`
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
		ReferenceID:   referenceID,
		ExternalID:    mtnResp.ExternalID,
		Status:        status,
		Amount:        parseFloat(mtnResp.Amount),
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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SubscriptionKey)

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

// GetResilienceStats returns resilience statistics
// Deprecated: Resilience is now handled by Temporal workflows
func (c *OrangeClient) GetResilienceStats() *models.ResilienceStats {
	// Return empty stats as resilience is handled by Temporal
	return &models.ResilienceStats{
		TotalRequests:      0,
		SuccessfulRequests: 0,
		FailedRequests:     0,
		LastReset:          time.Now(),
	}
}

// ResetResilienceStats resets resilience statistics
// Deprecated: Resilience is now handled by Temporal workflows
func (c *OrangeClient) ResetResilienceStats() {
	// No-op as resilience is handled by Temporal
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
			"value":    request.Amount,
			"unit":     request.Currency,
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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
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
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
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


