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

// CrudClient defines the interface for internal CRUD service client
type CrudClient interface {
	ClientInterface

	// CreateTransaction creates a new transaction in the CRUD service
	CreateTransaction(ctx context.Context, transaction *Transaction) (*TransactionResponse, error)

	// UpdateWalletBalance updates a wallet's balance in the CRUD service
	UpdateWalletBalance(ctx context.Context, walletID string, amount float64, currency string) (*WalletResponse, error)

	// GetWallet retrieves wallet information by user ID and currency
	GetWallet(ctx context.Context, userID string, currency string) (*WalletResponse, error)
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

// MobileMoneyConfig holds configuration for Mobile Money clients
type MobileMoneyConfig struct {
	Provider        string        // "mtn" or "orange"
	BaseURL         string        `yaml:"base_url"`
	APIKey          string        `yaml:"primary_key,omitempty"` // Pour MTN c'est primary_key, pour Orange c'est client_id
	APISecret       string        `yaml:"secondary_key,omitempty"` // Pour MTN c'est secondary_key, pour Orange c'est client_secret
	SubscriptionKey string        `yaml:"authorization,omitempty"` // Utilisé par Orange
	Timeout         time.Duration `yaml:"timeout"`
	MaxRetries      int           `yaml:"rate_limit"`
}

// Validate validates the Mobile Money configuration
func (c *MobileMoneyConfig) Validate() error {
	if c.Provider != "mtn" && c.Provider != "orange" && c.Provider != "" {
		return NewClientError("INVALID_PROVIDER", "provider must be 'mtn' or 'orange'", false)
	}
	if c.BaseURL == "" {
		return NewClientError("INVALID_CONFIG", "base URL is required", false)
	}
	if c.APIKey == "" {
		return NewClientError("INVALID_CONFIG", "API key is required", false)
	}
	if c.APISecret == "" {
		return NewClientError("INVALID_CONFIG", "API secret is required", false)
	}
	if c.SubscriptionKey == "" {
		return NewClientError("INVALID_CONFIG", "subscription key is required", false)
	}
	if c.Timeout <= 0 {
		return NewClientError("INVALID_CONFIG", "timeout must be positive", false)
	}
	if c.MaxRetries < 0 {
		return NewClientError("INVALID_CONFIG", "max retries cannot be negative", false)
	}
	return nil
}

// MobileMoneyError represents an error response from Mobile Money APIs
type MobileMoneyError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e MobileMoneyError) Error() string {
	return e.Message
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypePurchase TransactionType = "PURCHASE"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

// Transaction represents a financial transaction in the system
type Transaction struct {
	ID          string            `json:"id,omitempty"`
	Type        TransactionType   `json:"type" validate:"required"`
	Status      TransactionStatus `json:"status" validate:"required"`
	UserID      string            `json:"user_id" validate:"required"`
	WalletID    string            `json:"wallet_id,omitempty"`
	Amount      float64           `json:"amount" validate:"required,gt=0"`
	Currency    string            `json:"currency" validate:"required"`
	Reference   string            `json:"reference,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}

// Validate validates the transaction
func (t *Transaction) Validate() error {
	if t.Type == "" {
		return NewClientError("INVALID_TYPE", "Transaction type is required", false)
	}
	if t.Status == "" {
		return NewClientError("INVALID_STATUS", "Transaction status is required", false)
	}
	if t.UserID == "" {
		return NewClientError("INVALID_USER", "User ID is required", false)
	}
	if t.Amount <= 0 {
		return NewClientError("INVALID_AMOUNT", "Amount must be greater than 0", false)
	}
	if t.Currency == "" {
		return NewClientError("INVALID_CURRENCY", "Currency is required", false)
	}
	return nil
}

// TransactionResponse represents a response from creating or updating a transaction
type TransactionResponse struct {
	TransactionID string            `json:"transaction_id"`
	UserID        string            `json:"user_id"`
	WalletID      string            `json:"wallet_id,omitempty"`
	Type          TransactionType   `json:"type"`
	Status        TransactionStatus `json:"status"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Reference     string            `json:"reference,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
	Success       bool              `json:"success"`
	Error         string            `json:"error,omitempty"`
}

// WalletResponse represents a response from getting or updating a wallet
type WalletResponse struct {
	WalletID   string    `json:"wallet_id"`
	UserID     string    `json:"user_id"`
	Balance    float64   `json:"balance"`
	Currency   string    `json:"currency"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	Success    bool      `json:"success"`
	Error      string    `json:"error,omitempty"`
}
