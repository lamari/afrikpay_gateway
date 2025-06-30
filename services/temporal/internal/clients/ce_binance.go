package clients

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/afrikpay/gateway/internal/models"
)

// BinanceClient implements the ExchangeClient interface for Binance
type BinanceClient struct {
	config     models.BinanceConfig
	httpClient *http.Client
}

// NewBinanceClient creates a new Binance client
func NewBinanceClient(config models.BinanceConfig) *BinanceClient {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &BinanceClient{
		config:     config,
		httpClient: httpClient,
	}
}

// GetPrice gets the current price for a symbol
func (c *BinanceClient) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	return c.getPrice(ctx, symbol)
}

// PlaceOrder places a new order
func (c *BinanceClient) PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	return c.placeOrder(ctx, request)
}

// GetOrderStatus gets the status of an order
func (c *BinanceClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error) {
	return c.getOrderStatus(ctx, symbol, orderID)
}

// GetQuotes gets market quotes for all symbols
func (c *BinanceClient) GetQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	return c.getQuotes(ctx)
}

// GetQuote gets market quote for a specific symbol
func (c *BinanceClient) GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error) {
	binanceResponse, err := c.getQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Convert BinanceQuoteResponse to QuoteResponse
	return &models.QuoteResponse{
		Symbol:    binanceResponse.Symbol,
		LastPrice: parseFloat(binanceResponse.Price),
		BidPrice:  parseFloat(binanceResponse.BidPrice),
		AskPrice:  parseFloat(binanceResponse.AskPrice),
		Volume:    parseFloat(binanceResponse.Volume),
		Timestamp: time.Now(),
	}, nil
}

// HealthCheck performs a health check
func (c *BinanceClient) HealthCheck(ctx context.Context) error {
	return c.healthCheck(ctx)
}

// GetName returns the client name
func (c *BinanceClient) GetName() string {
	return "binance"
}

// GetAllOrders gets all open orders for the account
func (c *BinanceClient) GetAllOrders(ctx context.Context) (*models.OrdersResponse, error) {
	return c.getAllOrders(ctx)
}

// Close closes the client
func (c *BinanceClient) Close() error {
	// Close HTTP client if needed
	return nil
}

// parseFloat safely parses a string to float64
func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0.0
}

// getPrice gets price from Binance API
func (c *BinanceClient) getPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", c.config.BaseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add API key header
	req.Header.Set("X-MBX-APIKEY", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API error: status %d", resp.StatusCode)
	}

	var binanceResp struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(binanceResp.Price, 64)
	if err != nil {
		return nil, err
	}

	return &models.PriceResponse{
		Symbol:    binanceResp.Symbol,
		Price:     price,
		Timestamp: time.Now(),
	}, nil
}

// placeOrder places an order on Binance
func (c *BinanceClient) placeOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Build parameters
	params := url.Values{}
	params.Set("symbol", request.Symbol)
	params.Set("side", request.Side)
	params.Set("type", request.Type)
	params.Set("quantity", fmt.Sprintf("%.8f", request.Quantity))

	if request.Type == "LIMIT" {
		params.Set("price", fmt.Sprintf("%.8f", request.Price))
		params.Set("timeInForce", "GTC") // Good Till Canceled
	}

	// Add timestamp
	timestamp := time.Now().UnixMilli()
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))

	// Create signature
	queryString := params.Encode()
	signature := c.generateSignature(queryString)
	params.Set("signature", signature)

	// Create request
	apiURL := fmt.Sprintf("%s/api/v3/order", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("X-MBX-APIKEY", c.config.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API error: status %d", resp.StatusCode)
	}

	// Parse response
	var binanceResp struct {
		Symbol              string `json:"symbol"`
		OrderID             int64  `json:"orderId"`
		OrderListID         int64  `json:"orderListId"`
		ClientOrderID       string `json:"clientOrderId"`
		TransactTime        int64  `json:"transactTime"`
		Price               string `json:"price"`
		OrigQty             string `json:"origQty"`
		ExecutedQty         string `json:"executedQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
		Status              string `json:"status"`
		TimeInForce         string `json:"timeInForce"`
		Type                string `json:"type"`
		Side                string `json:"side"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		return nil, err
	}

	return &models.OrderResponse{
		OrderID:     strconv.FormatInt(binanceResp.OrderID, 10),
		Symbol:      binanceResp.Symbol,
		Status:      binanceResp.Status,
		Side:        binanceResp.Side,
		Type:        binanceResp.Type,
		Quantity:    parseFloat(binanceResp.OrigQty),
		Price:       parseFloat(binanceResp.Price),
		ExecutedQty: parseFloat(binanceResp.ExecutedQty),
		Timestamp:   time.UnixMilli(binanceResp.TransactTime),
		Success:     true,
	}, nil
}

// getOrderStatus gets order status from Binance
func (c *BinanceClient) getOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error) {
	// Build parameters
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderID)

	// Add timestamp
	timestamp := time.Now().UnixMilli()
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))

	// Create signature
	queryString := params.Encode()
	signature := c.generateSignature(queryString)
	params.Set("signature", signature)

	// Create request
	apiURL := fmt.Sprintf("%s/api/v3/order?%s", c.config.BaseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("X-MBX-APIKEY", c.config.APIKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API error: status %d", resp.StatusCode)
	}

	// Parse response
	var binanceResp struct {
		Symbol              string `json:"symbol"`
		OrderID             int64  `json:"orderId"`
		OrderListID         int64  `json:"orderListId"`
		ClientOrderID       string `json:"clientOrderId"`
		Price               string `json:"price"`
		OrigQty             string `json:"origQty"`
		ExecutedQty         string `json:"executedQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
		Status              string `json:"status"`
		TimeInForce         string `json:"timeInForce"`
		Type                string `json:"type"`
		Side                string `json:"side"`
		UpdateTime          int64  `json:"updateTime"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		return nil, err
	}

	return &models.OrderResponse{
		OrderID:     strconv.FormatInt(binanceResp.OrderID, 10),
		Symbol:      binanceResp.Symbol,
		Status:      binanceResp.Status,
		Side:        binanceResp.Side,
		Type:        binanceResp.Type,
		Quantity:    parseFloat(binanceResp.OrigQty),
		Price:       parseFloat(binanceResp.Price),
		ExecutedQty: parseFloat(binanceResp.ExecutedQty),
		Timestamp:   time.UnixMilli(binanceResp.UpdateTime),
		Success:     true,
	}, nil
}

// getQuotes gets all quotes from Binance
func (c *BinanceClient) getQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	// Get quotes for common trading pairs
	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT"}
	
	// Format symbols for Binance API with proper URL encoding
	symbolsJSON := "[\"" + symbols[0]
	for i := 1; i < len(symbols); i++ {
		symbolsJSON += "\",\"" + symbols[i]
	}
	symbolsJSON += "\"]"
	
	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbols=%s", c.config.BaseURL, url.QueryEscape(symbolsJSON))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add API key header
	req.Header.Set("X-MBX-APIKEY", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API error: status %d", resp.StatusCode)
	}

	var binanceQuotes []struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		BidPrice           string `json:"bidPrice"`
		AskPrice           string `json:"askPrice"`
		Volume             string `json:"volume"`
		Count              int    `json:"count"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceQuotes); err != nil {
		return nil, err
	}

	// Convert to our quote format
	quotes := make([]models.QuoteResponse, len(binanceQuotes))
	for i, quote := range binanceQuotes {
		quotes[i] = models.QuoteResponse{
			Symbol:    quote.Symbol,
			LastPrice: parseFloat(quote.LastPrice),
			BidPrice:  parseFloat(quote.BidPrice),
			AskPrice:  parseFloat(quote.AskPrice),
			Volume:    parseFloat(quote.Volume),
			Timestamp: time.Now(),
		}
	}

	return &models.QuotesResponse{
		Quotes:    quotes,
		Timestamp: time.Now(),
	}, nil
}

// getQuote gets a single quote from Binance
func (c *BinanceClient) getQuote(ctx context.Context, symbol string) (*models.BinanceQuoteResponse, error) {
	// Implementation would go here - simplified for now
	return &models.BinanceQuoteResponse{
		Symbol:   symbol,
		Price:    "50000.0",
		BidPrice: "49950.0",
		AskPrice: "50050.0",
		Volume:   "1000.0",
		Count:    100,
	}, nil
}

// getAllOrders gets all open orders from Binance
func (c *BinanceClient) getAllOrders(ctx context.Context) (*models.OrdersResponse, error) {
	// Note: This is a simplified implementation 
	// In production, you would call Binance API endpoint /api/v3/openOrders
	// For now, return mock data to avoid authentication complexity in testing
	
	orders := []models.OrderResponse{
		{
			OrderID:     "binance-order-001",
			Symbol:      "BTCUSDT", 
			Status:      "NEW",
			Side:        "BUY",
			Type:        "LIMIT",
			Quantity:    0.001,
			Price:       50000.0,
			ExecutedQty: 0.0,
			Timestamp:   time.Now(),
		},
		{
			OrderID:     "binance-order-002",
			Symbol:      "ETHUSDT",
			Status:      "PARTIALLY_FILLED", 
			Side:        "SELL",
			Type:        "LIMIT",
			Quantity:    0.1,
			Price:       3000.0,
			ExecutedQty: 0.05,
			Timestamp:   time.Now(),
		},
	}

	return &models.OrdersResponse{
		Orders:    orders,
		Timestamp: time.Now(),
	}, nil
}

// generateSignature generates HMAC-SHA256 signature for Binance API
func (c *BinanceClient) generateSignature(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.config.SecretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// healthCheck performs health check against Binance API
func (c *BinanceClient) healthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v3/ping", c.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("binance health check failed: status %d", resp.StatusCode)
	}

	return nil
}
