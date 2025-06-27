package models

import (
	"context"
	"time"
)

// ClientInterface defines the common interface for all API clients
type ClientInterface interface {
	// Health check
	HealthCheck(ctx context.Context) error
	
	// Get client name for logging
	GetName() string
	
	// Get resilience statistics
	GetResilienceStats() *ResilienceStats
	
	// Reset resilience statistics
	ResetResilienceStats()
	
	// Close client and cleanup resources
	Close() error
}

// ExchangeClient defines the interface for cryptocurrency exchange clients
type ExchangeClient interface {
	ClientInterface
	
	// Get current price for a symbol
	GetPrice(ctx context.Context, symbol string) (*PriceResponse, error)
	
	// Place an order
	PlaceOrder(ctx context.Context, order *OrderRequest) (*OrderResponse, error)
	
	// Get order status
	GetOrderStatus(ctx context.Context, symbol string, orderID string) (*OrderResponse, error)
	
	// Get quotes for all symbols
	GetQuotes(ctx context.Context) (*QuotesResponse, error)
	
	// Get quote for a specific symbol
	GetQuote(ctx context.Context, symbol string) (*QuoteResponse, error)
}

// MobileMoneyClient defines the interface for mobile money clients
type MobileMoneyClient interface {
	ClientInterface
	
	// Initiate payment
	InitiatePayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)
	
	// Get payment status
	GetPaymentStatus(ctx context.Context, referenceID string) (*PaymentStatusResponse, error)
}

// PriceResponse represents a generic price response
type PriceResponse struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// OrderRequest represents a generic order request
type OrderRequest struct {
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"` // BUY or SELL
	Type        string  `json:"type"` // MARKET or LIMIT
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price,omitempty"`
	TimeInForce string  `json:"timeInForce,omitempty"`
}

// Validate validates the order request
func (o *OrderRequest) Validate() error {
	if o.Symbol == "" {
		return NewClientError("INVALID_SYMBOL", "Symbol is required", false)
	}
	if o.Side != string(OrderSideBuy) && o.Side != string(OrderSideSell) {
		return NewClientError("INVALID_SIDE", "Side must be BUY or SELL", false)
	}
	if o.Type != string(OrderTypeMarket) && o.Type != string(OrderTypeLimit) {
		return NewClientError("INVALID_TYPE", "Type must be MARKET or LIMIT", false)
	}
	if o.Quantity <= 0 {
		return NewClientError("INVALID_QUANTITY", "Quantity must be greater than 0", false)
	}
	if o.Type == string(OrderTypeLimit) && o.Price <= 0 {
		return NewClientError("INVALID_PRICE", "Price is required for LIMIT orders", false)
	}
	return nil
}

// OrderResponse represents a generic order response
type OrderResponse struct {
	OrderID       string    `json:"orderId"`
	Symbol        string    `json:"symbol"`
	Status        string    `json:"status"`
	Side          string    `json:"side"`
	Type          string    `json:"type"`
	Quantity      float64   `json:"quantity"`
	Price         float64   `json:"price"`
	ExecutedQty   float64   `json:"executedQty"`
	Timestamp     time.Time `json:"timestamp"`
	ClientOrderID string    `json:"clientOrderId,omitempty"`
	Success       bool      `json:"success"`
	Error         string    `json:"error,omitempty"`
	AvgPrice      float64   `json:"avgPrice,omitempty"`
	Fee           float64   `json:"fee,omitempty"`
	FeeCurrency   string    `json:"feeCurrency,omitempty"`
}

// PaymentRequest represents a generic payment request
type PaymentRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	ExternalID  string            `json:"externalId"`
	PhoneNumber string            `json:"phoneNumber"`
	Description string            `json:"description,omitempty"`
	CallbackURL string            `json:"callbackUrl,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PaymentResponse represents a generic payment response
type PaymentResponse struct {
	ReferenceID string        `json:"referenceId"`
	Status      PaymentStatus `json:"status"`
	PaymentURL  string        `json:"paymentUrl,omitempty"`
	Message     string        `json:"message,omitempty"`
}

// PaymentStatusResponse represents a generic payment status response
type PaymentStatusResponse struct {
	ReferenceID   string        `json:"referenceId"`
	ExternalID    string        `json:"externalId"`
	Status        PaymentStatus `json:"status"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PhoneNumber   string        `json:"phoneNumber"`
	TransactionID string        `json:"transactionId,omitempty"`
	Message       string        `json:"message,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
}

// ClientError represents a generic client error
type ClientError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Retryable  bool   `json:"retryable"`
}

func (e ClientError) Error() string {
	return e.Message
}

// IsRetryable returns true if the error is retryable
func (e ClientError) IsRetryable() bool {
	return e.Retryable
}

// NewClientError creates a new ClientError
func NewClientError(code, message string, retryable bool) *ClientError {
	return &ClientError{
		Code:      code,
		Message:   message,
		Retryable: retryable,
	}
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// String returns the string representation of OrderStatus
func (o OrderStatus) String() string {
	return string(o)
}

// OrderSide represents the side of an order
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// String returns the string representation of OrderSide
func (o OrderSide) String() string {
	return string(o)
}

// OrderType represents the type of an order
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

// String returns the string representation of OrderType
func (o OrderType) String() string {
	return string(o)
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCanceled   PaymentStatus = "CANCELED"
	PaymentStatusExpired    PaymentStatus = "EXPIRED"
)

// String returns the string representation of PaymentStatus
func (p PaymentStatus) String() string {
	return string(p)
}

// DepositRequest represents a deposit request for mobile money
type DepositRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	PhoneNumber string            `json:"phoneNumber"`
	Description string            `json:"description,omitempty"`
	ExternalID  string            `json:"externalId"`
	CallbackURL string            `json:"callbackUrl,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Validate validates the deposit request
func (d *DepositRequest) Validate() error {
	if d.Amount <= 0 {
		return NewClientError("INVALID_AMOUNT", "Amount must be greater than 0", false)
	}
	if d.Currency == "" {
		return NewClientError("INVALID_CURRENCY", "Currency is required", false)
	}
	if d.PhoneNumber == "" {
		return NewClientError("INVALID_PHONE", "Phone number is required", false)
	}
	if d.ExternalID == "" {
		return NewClientError("INVALID_EXTERNAL_ID", "External ID is required", false)
	}
	return nil
}

// DepositResponse represents a deposit response
type DepositResponse struct {
	TransactionID string        `json:"transactionId"`
	ReferenceID   string        `json:"referenceId"`
	Status        PaymentStatus `json:"status"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PhoneNumber   string        `json:"phoneNumber"`
	Message       string        `json:"message,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
}

// QuoteResponse represents a quote response for multiple symbols
type QuoteResponse struct {
	Symbol    string    `json:"symbol"`
	BidPrice  float64   `json:"bidPrice"`
	AskPrice  float64   `json:"askPrice"`
	LastPrice float64   `json:"lastPrice"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// QuotesResponse represents a response containing multiple quotes
type QuotesResponse struct {
	Quotes    []QuoteResponse `json:"quotes"`
	Timestamp time.Time       `json:"timestamp"`
}

// OrderStatusResponse represents the response for order status queries
type OrderStatusResponse struct {
	OrderID       string    `json:"orderId"`
	Symbol        string    `json:"symbol"`
	Status        string    `json:"status"`
	Side          string    `json:"side"`
	Type          string    `json:"type"`
	Quantity      float64   `json:"quantity"`
	Price         float64   `json:"price"`
	ExecutedQty   float64   `json:"executedQty"`
	Timestamp     time.Time `json:"timestamp"`
	ClientOrderID string    `json:"clientOrderId,omitempty"`
}

// TransactionStatusResponse represents the response for transaction status queries
type TransactionStatusResponse struct {
	TransactionID string        `json:"transaction_id"`
	ReferenceID   string        `json:"reference_id"`
	Status        PaymentStatus `json:"status"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Timestamp     time.Time     `json:"timestamp"`
	Message       string        `json:"message,omitempty"`
}

// ResilienceStats represents resilience statistics
type ResilienceStats struct {
	TotalRequests       int64     `json:"total_requests"`
	SuccessfulRequests  int64     `json:"successful_requests"`
	FailedRequests      int64     `json:"failed_requests"`
	CircuitBreakerTrips int64     `json:"circuit_breaker_trips"`
	RetryAttempts       int64     `json:"retry_attempts"`
	LastReset           time.Time `json:"last_reset"`
	CircuitBreakerState string    `json:"circuit_breaker_state"`
}

// Reset resets all statistics
func (r *ResilienceStats) Reset() {
	r.TotalRequests = 0
	r.SuccessfulRequests = 0
	r.FailedRequests = 0
	r.CircuitBreakerTrips = 0
	r.RetryAttempts = 0
	r.LastReset = time.Now()
}
