package clients

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/afrikpay/gateway/internal/models"
)

// BitgetClient represents a client for Bitget exchange API
type BitgetClient struct {
	config     *models.BitgetConfig
	httpClient *http.Client
}

// NewBitgetClient creates a new Bitget client
func NewBitgetClient(config *models.BitgetConfig) (*BitgetClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &BitgetClient{
		config:     config,
		httpClient: httpClient,
	}, nil
}

// GetPrice gets the current price for a symbol
func (c *BitgetClient) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	request := &models.BitgetPriceRequest{
		Symbol: symbol,
	}

	response, err := c.getPrice(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.ToPriceResponse(), nil
}

// PlaceOrder places a new order
func (c *BitgetClient) PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	// Convert common OrderRequest to BitgetOrderRequest
	bitgetRequest := &models.BitgetOrderRequest{
		Symbol:    request.Symbol,
		Side:      string(request.Side),
		OrderType: string(request.Type),
		Size:      fmt.Sprintf("%.8f", request.Quantity),
		Price:     fmt.Sprintf("%.8f", request.Price),
	}

	response, err := c.placeOrder(ctx, bitgetRequest)

	if err != nil {
		return nil, err
	}

	return response.ToOrderResponse(), nil
}

// GetOrderStatus gets the status of an order
func (c *BitgetClient) GetOrderStatus(ctx context.Context, symbol, orderID string) (*models.OrderResponse, error) {
	request := &models.BitgetOrderStatusRequest{
		Symbol:  symbol,
		OrderId: orderID,
	}

	response, err := c.getOrderStatus(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.ToOrderResponse(), nil
}

// GetQuotes gets market quotes for all symbols
func (c *BitgetClient) GetQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	bitgetResponse, err := c.getQuotes(ctx)

	if err != nil {
		return nil, err
	}

	// Convert BitgetQuotesResponse to common QuotesResponse
	quotes := make([]models.QuoteResponse, len(bitgetResponse.Data))
	for i, quote := range bitgetResponse.Data {
		// Convert string prices to float64
		price, _ := strconv.ParseFloat(quote.Price, 64)
		volume, _ := strconv.ParseFloat(quote.Volume, 64)

		quotes[i] = models.QuoteResponse{
			Symbol:    quote.Symbol,
			LastPrice: price,
			BidPrice:  price, // Bitget doesn't provide separate bid/ask, use price
			AskPrice:  price, // Bitget doesn't provide separate bid/ask, use price
			Volume:    volume,
			Timestamp: time.Now(),
		}
	}

	return &models.QuotesResponse{
		Quotes: quotes,
	}, nil
}

// GetQuote gets market quote for a specific symbol
func (c *BitgetClient) GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error) {
	bitgetResponse, err := c.getQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Convert BitgetQuoteResponse to QuoteResponse
	response := &models.QuoteResponse{
		Symbol:    bitgetResponse.Data.Symbol,
		LastPrice: parseFloat(bitgetResponse.Data.Price),
		BidPrice:  0.0, // Not available in this response
		AskPrice:  0.0, // Not available in this response
		Volume:    parseFloat(bitgetResponse.Data.Volume),
		Timestamp: time.Now(),
	}

	return response, nil
}

// getPrice makes the actual HTTP request to get price
func (c *BitgetClient) getPrice(ctx context.Context, request *models.BitgetPriceRequest) (*models.BitgetPriceResponse, error) {
	endpoint := "/api/spot/v1/market/ticker"
	params := map[string]string{
		"symbol": request.Symbol,
	}

	var response models.BitgetPriceResponse
	err := c.makeRequest(ctx, "GET", endpoint, params, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// placeOrder makes the actual HTTP request to place an order
func (c *BitgetClient) placeOrder(ctx context.Context, request *models.BitgetOrderRequest) (*models.BitgetOrderResponse, error) {
	endpoint := "/api/spot/v1/trade/orders"

	var response models.BitgetOrderResponse
	err := c.makeRequest(ctx, "POST", endpoint, nil, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getOrderStatus makes the actual HTTP request to get order status
func (c *BitgetClient) getOrderStatus(ctx context.Context, request *models.BitgetOrderStatusRequest) (*models.BitgetOrderStatusResponse, error) {
	endpoint := "/api/spot/v1/trade/orderInfo"
	params := map[string]string{
		"symbol":  request.Symbol,
		"orderId": request.OrderId,
	}

	var response models.BitgetOrderStatusResponse
	err := c.makeRequest(ctx, "GET", endpoint, params, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getQuotes makes the actual HTTP request to get all market quotes
func (c *BitgetClient) getQuotes(ctx context.Context) (*models.BitgetQuotesResponse, error) {
	endpoint := "/api/spot/v1/market/tickers"

	var response models.BitgetQuotesResponse
	err := c.makeRequest(ctx, "GET", endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// getQuote makes the actual HTTP request to get a specific market quote
func (c *BitgetClient) getQuote(ctx context.Context, symbol string) (*models.BitgetQuoteResponse, error) {
	endpoint := "/api/spot/v1/market/ticker"
	params := map[string]string{
		"symbol": symbol,
	}

	var response models.BitgetQuoteResponse
	err := c.makeRequest(ctx, "GET", endpoint, params, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// makeRequest makes an HTTP request with proper authentication
func (c *BitgetClient) makeRequest(ctx context.Context, method, endpoint string, params map[string]string, body interface{}, result interface{}) error {
	// Build URL
	url := c.config.BaseURL + endpoint
	if len(params) > 0 {
		url += "?"
		for key, value := range params {
			url += key + "=" + value + "&"
		}
		url = url[:len(url)-1] // Remove trailing &
	}

	// Prepare request body
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return models.NewClientError("MARSHAL_ERROR", fmt.Sprintf("failed to marshal request body: %v", err), false)
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.NewClientError("REQUEST_ERROR", fmt.Sprintf("failed to create request: %v", err), true)
	}

	// Add authentication headers
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := c.generateSignature(method, endpoint, timestamp, string(reqBody))

	req.Header.Set("ACCESS-KEY", c.config.APIKey)
	req.Header.Set("ACCESS-SIGN", signature)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("ACCESS-PASSPHRASE", c.config.Passphrase)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.NewClientError("HTTP_ERROR", fmt.Sprintf("HTTP request failed: %v", err), true)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.NewClientError("READ_ERROR", fmt.Sprintf("failed to read response body: %v", err), true)
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		var errorResp models.BitgetErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			retryable := resp.StatusCode >= 500 || resp.StatusCode == 429
			return models.NewClientError("API_ERROR", fmt.Sprintf("API error: %s - %s", errorResp.Code, errorResp.Message), retryable)
		}
		retryable := resp.StatusCode >= 500 || resp.StatusCode == 429
		return models.NewClientError("HTTP_ERROR", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody)), retryable)
	}

	// Parse response
	if err := json.Unmarshal(respBody, result); err != nil {
		return models.NewClientError("UNMARSHAL_ERROR", fmt.Sprintf("failed to unmarshal response: %v", err), false)
	}

	return nil
}

// generateSignature generates the signature for Bitget API authentication
func (c *BitgetClient) generateSignature(method, endpoint, timestamp, body string) string {
	// Create the message to sign: timestamp + method + endpoint + body
	message := timestamp + method + endpoint + body

	// Create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(c.config.SecretKey))
	h.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// HealthCheck performs a health check
func (c *BitgetClient) HealthCheck(ctx context.Context) error {
	// Simple health check - for now, just return nil - can be enhanced with actual health check
	return nil
}

// Close closes the client and cleans up resources
func (c *BitgetClient) Close() error {
	// Close HTTP client if it has a Close method
	if closer, ok := c.httpClient.Transport.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
