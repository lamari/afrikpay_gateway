package models

import (
	"strconv"
	"time"
)

// BitgetConfig holds configuration for Bitget API client
type BitgetConfig struct {
	BaseURL    string        `json:"base_url" validate:"required,url"`
	APIKey     string        `json:"api_key" validate:"required"`
	SecretKey  string        `json:"secret_key" validate:"required"`
	Passphrase string        `json:"passphrase" validate:"required"`
	Timeout    time.Duration `json:"timeout" validate:"required"`
	MaxRetries int           `json:"max_retries" validate:"min=0,max=10"`
}

// BitgetPriceRequest represents a request to get current price from Bitget
type BitgetPriceRequest struct {
	Symbol string `json:"symbol" validate:"required"`
}

// BitgetPriceResponse represents the response from Bitget price API
type BitgetPriceResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	} `json:"data"`
	RequestTime int64 `json:"requestTime"`
}

// BitgetOrderRequest represents a request to place an order on Bitget
type BitgetOrderRequest struct {
	Symbol      string `json:"symbol" validate:"required"`
	Side        string `json:"side" validate:"required,oneof=buy sell"`
	OrderType   string `json:"orderType" validate:"required,oneof=market limit"`
	Size        string `json:"size" validate:"required"`
	Price       string `json:"price,omitempty"`
	ClientOid   string `json:"clientOid,omitempty"`
	TimeInForce string `json:"timeInForce,omitempty"`
}

// BitgetOrderResponse represents the response from Bitget order API
type BitgetOrderResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		OrderId       string `json:"orderId"`
		ClientOid     string `json:"clientOid"`
		Symbol        string `json:"symbol"`
		Side          string `json:"side"`
		OrderType     string `json:"orderType"`
		Size          string `json:"size"`
		Price         string `json:"price"`
		Status        string `json:"status"`
		FilledSize    string `json:"filledSize"`
		FilledAmount  string `json:"filledAmount"`
		CreateTime    int64  `json:"createTime"`
		UpdateTime    int64  `json:"updateTime"`
	} `json:"data"`
	RequestTime int64 `json:"requestTime"`
}

// BitgetOrderStatusRequest represents a request to get order status from Bitget
type BitgetOrderStatusRequest struct {
	Symbol  string `json:"symbol" validate:"required"`
	OrderId string `json:"orderId" validate:"required"`
}

// BitgetOrderStatusResponse represents the response from Bitget order status API
type BitgetOrderStatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		OrderId       string `json:"orderId"`
		ClientOid     string `json:"clientOid"`
		Symbol        string `json:"symbol"`
		Side          string `json:"side"`
		OrderType     string `json:"orderType"`
		Size          string `json:"size"`
		Price         string `json:"price"`
		Status        string `json:"status"`
		FilledSize    string `json:"filledSize"`
		FilledAmount  string `json:"filledAmount"`
		AvgPrice      string `json:"avgPrice"`
		Fee           string `json:"fee"`
		FeeCurrency   string `json:"feeCurrency"`
		CreateTime    int64  `json:"createTime"`
		UpdateTime    int64  `json:"updateTime"`
	} `json:"data"`
	RequestTime int64 `json:"requestTime"`
}

// BitgetQuotesResponse represents the response from Bitget quotes API
type BitgetQuotesResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    []struct {
		Symbol   string `json:"symbol"`
		Price    string `json:"price"`
		Change   string `json:"change"`
		ChangeP  string `json:"changeP"`
		High24h  string `json:"high24h"`
		Low24h   string `json:"low24h"`
		Volume   string `json:"volume"`
		Turnover string `json:"turnover"`
	} `json:"data"`
	RequestTime int64 `json:"requestTime"`
}

// BitgetQuoteResponse represents the response from Bitget single quote API
type BitgetQuoteResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Symbol   string `json:"symbol"`
		Price    string `json:"price"`
		Change   string `json:"change"`
		ChangeP  string `json:"changeP"`
		High24h  string `json:"high24h"`
		Low24h   string `json:"low24h"`
		Volume   string `json:"volume"`
		Turnover string `json:"turnover"`
	} `json:"data"`
	RequestTime int64 `json:"requestTime"`
}

// BitgetErrorResponse represents an error response from Bitget API
type BitgetErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// Bitget order status constants
const (
	BitgetOrderStatusNew             = "new"
	BitgetOrderStatusPartiallyFilled = "partially_filled"
	BitgetOrderStatusFilled          = "filled"
	BitgetOrderStatusCanceled        = "canceled"
	BitgetOrderStatusRejected        = "rejected"
)

// Bitget order side constants
const (
	BitgetOrderSideBuy  = "buy"
	BitgetOrderSideSell = "sell"
)

// Bitget order type constants
const (
	BitgetOrderTypeMarket = "market"
	BitgetOrderTypeLimit  = "limit"
)

// Bitget time in force constants
const (
	BitgetTimeInForceGTC = "GTC" // Good Till Canceled
	BitgetTimeInForceIOC = "IOC" // Immediate Or Cancel
	BitgetTimeInForceFOK = "FOK" // Fill Or Kill
)

// Validate validates the BitgetConfig
func (c *BitgetConfig) Validate() error {
	if c.BaseURL == "" {
		return NewClientError("INVALID_CONFIG", "base URL is required", false)
	}
	if c.APIKey == "" {
		return NewClientError("INVALID_CONFIG", "API key is required", false)
	}
	if c.SecretKey == "" {
		return NewClientError("INVALID_CONFIG", "secret key is required", false)
	}
	if c.Passphrase == "" {
		return NewClientError("INVALID_CONFIG", "passphrase is required", false)
	}
	if c.Timeout <= 0 {
		return NewClientError("INVALID_CONFIG", "timeout must be positive", false)
	}
	if c.MaxRetries < 0 {
		return NewClientError("INVALID_CONFIG", "max retries cannot be negative", false)
	}
	return nil
}

// IsSuccess checks if the Bitget response indicates success
func (r *BitgetPriceResponse) IsSuccess() bool {
	return r.Code == "00000"
}

// IsSuccess checks if the Bitget response indicates success
func (r *BitgetOrderResponse) IsSuccess() bool {
	return r.Code == "00000"
}

// IsSuccess checks if the Bitget response indicates success
func (r *BitgetOrderStatusResponse) IsSuccess() bool {
	return r.Code == "00000"
}

// IsSuccess checks if the Bitget response indicates success
func (r *BitgetQuotesResponse) IsSuccess() bool {
	return r.Code == "00000"
}

// IsSuccess checks if the Bitget response indicates success
func (r *BitgetQuoteResponse) IsSuccess() bool {
	return r.Code == "00000"
}

// ToOrderResponse converts BitgetOrderResponse to common OrderResponse
func (r *BitgetOrderResponse) ToOrderResponse() *OrderResponse {
	if !r.IsSuccess() {
		// Return an empty OrderResponse for failed requests
		// The error should be handled at the client level
		return &OrderResponse{}
	}

	// Convert string values to appropriate types
	quantity, _ := strconv.ParseFloat(r.Data.Size, 64)
	price, _ := strconv.ParseFloat(r.Data.Price, 64)
	executedQty, _ := strconv.ParseFloat(r.Data.FilledSize, 64)

	return &OrderResponse{
		OrderID:       r.Data.OrderId,
		Symbol:        r.Data.Symbol,
		Side:          r.Data.Side,
		Type:          r.Data.OrderType,
		Quantity:      quantity,
		Price:         price,
		Status:        r.Data.Status,
		ExecutedQty:   executedQty,
		Timestamp:     time.Unix(r.Data.CreateTime/1000, 0),
		ClientOrderID: r.Data.ClientOid,
	}
}

// ToOrderResponse converts BitgetOrderStatusResponse to common OrderResponse
func (r *BitgetOrderStatusResponse) ToOrderResponse() *OrderResponse {
	if !r.IsSuccess() {
		// Return an empty OrderResponse for failed requests
		// The error should be handled at the client level
		return &OrderResponse{}
	}

	// Convert string values to appropriate types
	quantity, _ := strconv.ParseFloat(r.Data.Size, 64)
	price, _ := strconv.ParseFloat(r.Data.Price, 64)
	executedQty, _ := strconv.ParseFloat(r.Data.FilledSize, 64)

	return &OrderResponse{
		OrderID:       r.Data.OrderId,
		Symbol:        r.Data.Symbol,
		Side:          r.Data.Side,
		Type:          r.Data.OrderType,
		Quantity:      quantity,
		Price:         price,
		Status:        r.Data.Status,
		ExecutedQty:   executedQty,
		Timestamp:     time.Unix(r.Data.CreateTime/1000, 0),
		ClientOrderID: r.Data.ClientOid,
	}
}

// ToPriceResponse converts BitgetPriceResponse to common PriceResponse
func (r *BitgetPriceResponse) ToPriceResponse() *PriceResponse {
	if !r.IsSuccess() {
		// Return an empty PriceResponse for failed requests
		// The error should be handled at the client level
		return &PriceResponse{}
	}

	// Convert string price to float64
	price, _ := strconv.ParseFloat(r.Data.Price, 64)

	return &PriceResponse{
		Symbol:    r.Data.Symbol,
		Price:     price,
		Timestamp: time.Unix(r.RequestTime/1000, 0),
	}
}

// ToQuotesResponse converts BitgetQuotesResponse to common QuotesResponse
func (r *BitgetQuotesResponse) ToQuotesResponse() *QuotesResponse {
	if !r.IsSuccess() {
		// Return an empty QuotesResponse for failed requests
		return &QuotesResponse{
			Quotes:    []QuoteResponse{},
			Timestamp: time.Now(),
		}
	}

	quotes := make([]QuoteResponse, len(r.Data))
	for i, data := range r.Data {
		// Convert string values to float64
		bidPrice, _ := strconv.ParseFloat(data.Price, 64)
		askPrice, _ := strconv.ParseFloat(data.Price, 64) // Bitget doesn't separate bid/ask in this response
		lastPrice, _ := strconv.ParseFloat(data.Price, 64)
		volume, _ := strconv.ParseFloat(data.Volume, 64)

		quotes[i] = QuoteResponse{
			Symbol:    data.Symbol,
			BidPrice:  bidPrice,
			AskPrice:  askPrice,
			LastPrice: lastPrice,
			Volume:    volume,
			Timestamp: time.Unix(r.RequestTime/1000, 0),
		}
	}

	return &QuotesResponse{
		Quotes:    quotes,
		Timestamp: time.Unix(r.RequestTime/1000, 0),
	}
}

// ToQuoteResponse converts BitgetQuoteResponse to common QuoteResponse
func (r *BitgetQuoteResponse) ToQuoteResponse() *QuoteResponse {
	if !r.IsSuccess() {
		// Return an empty QuoteResponse for failed requests
		return &QuoteResponse{}
	}

	// Convert string values to float64
	bidPrice, _ := strconv.ParseFloat(r.Data.Price, 64)
	askPrice, _ := strconv.ParseFloat(r.Data.Price, 64) // Bitget doesn't separate bid/ask in this response
	lastPrice, _ := strconv.ParseFloat(r.Data.Price, 64)
	volume, _ := strconv.ParseFloat(r.Data.Volume, 64)

	return &QuoteResponse{
		Symbol:    r.Data.Symbol,
		BidPrice:  bidPrice,
		AskPrice:  askPrice,
		LastPrice: lastPrice,
		Volume:    volume,
		Timestamp: time.Unix(r.RequestTime/1000, 0),
	}
}
