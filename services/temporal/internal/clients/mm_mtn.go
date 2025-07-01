package clients

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

// GenerateTimestampUUID generates a UUID with timestamp information embedded
func GenerateTimestampUUID() string {
	// Utiliser le timestamp Unix en nanosecondes
	timestamp := time.Now().UnixNano()

	// Générer des bytes aléatoires
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// Combiner timestamp + random pour créer un UUID valide
	uuidBytes := make([]byte, 16)

	// 8 premiers bytes: timestamp (little endian)
	for i := 0; i < 8; i++ {
		uuidBytes[i] = byte(timestamp >> (i * 8))
	}

	// 8 derniers bytes: aléatoires
	copy(uuidBytes[8:], randomBytes)

	// Forcer le format UUID v4
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40 // Version 4
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant bits

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuidBytes[0:4],
		uuidBytes[4:6],
		uuidBytes[6:8],
		uuidBytes[8:10],
		uuidBytes[10:16])
}

// CreateUser creates a new API user in MTN MoMo and stores the reference ID
func (c *MTNClient) CreateUser(ctx context.Context, providerCallbackHost string) (string, error) {
	referenceID := GenerateTimestampUUID()
	url := fmt.Sprintf("%s/v1_0/apiuser", c.config.BaseURL)
	payload := map[string]string{"providerCallbackHost": providerCallbackHost}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("MTN CreateUser API error: status %d", resp.StatusCode)
	}

	return referenceID, nil
}

// CreateApiKey generates an API key for the MTN API user and stores it
func (c *MTNClient) CreateApiKey(ctx context.Context, referenceID string) (string, error) {
	if referenceID == "" {
		return "", fmt.Errorf("MTN reference ID is not set")
	}
	url := fmt.Sprintf("%s/v1_0/apiuser/%s/apikey", c.config.BaseURL, referenceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)
	req.Header.Set("X-Reference-Id", referenceID)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("MTN CreateApiKey API error: status %d", resp.StatusCode)
	}
	var res struct {
		ApiKey string `json:"apiKey"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.ApiKey, nil
}

// GetAccessToken retrieves an access token using the API key and stores it
func (c *MTNClient) GetAccessToken(ctx context.Context, referenceID, apiKey string) (string, error) {
	if referenceID == "" || apiKey == "" {
		return "", fmt.Errorf("MTN reference ID or API key is not set")
	}
	url := fmt.Sprintf("%s/collection/token/", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)
	credential := base64.StdEncoding.EncodeToString([]byte(referenceID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+credential)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("MTN GetAccessToken API error: status %d", resp.StatusCode)
	}
	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

// CreatePaymentRequest sends a payment request and returns the MTN response
func (c *MTNClient) CreatePaymentRequest(ctx context.Context, referenceID, accessToken string, request *models.MTNPaymentRequest) (*models.MTNPaymentResponse, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("MTN access token is not set")
	}
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", c.config.BaseURL)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Debug logging
	fmt.Printf("[DEBUG] CreatePaymentRequest URL: %s\n", url)
	fmt.Printf("[DEBUG] CreatePaymentRequest ReferenceID: %s\n", referenceID)
	fmt.Printf("[DEBUG] CreatePaymentRequest AccessToken: %s\n", accessToken[:10]+"...")
	fmt.Printf("[DEBUG] CreatePaymentRequest Payload: %s\n", string(jsonData))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.config.SecondaryKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Debug: read response body for error details
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("[DEBUG] CreatePaymentRequest Response Status: %d\n", resp.StatusCode)
	fmt.Printf("[DEBUG] CreatePaymentRequest Response Body: %s\n", string(respBody))

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("MTN CreatePaymentRequest API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Si nous avons un code 2xx, c'est un succès même avec un corps vide
	var res models.MTNPaymentResponse
	
	// Si le corps n'est pas vide, on essaie de le décoder
	if len(respBody) > 0 {
		// Reset body for JSON decoding
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		
		// Tenter de décoder mais ne pas échouer si erreur EOF
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil && err != io.EOF {
			// On ne retourne une erreur que si c'est autre chose que EOF
			return nil, fmt.Errorf("erreur de décodage JSON: %w", err)
		}
	}
	
	// Si le corps est vide ou ne contenait rien d'important, on remplit au moins les champs essentiels
	if res.Status == "" {
		// Pour une réponse avec code 202, on considère que le statut est PENDING
		res.Status = string(models.PaymentStatusPending)
		res.ReferenceID = referenceID
		res.Reason = "Payment request accepted by MTN API"
	}
	
	return &res, nil
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
