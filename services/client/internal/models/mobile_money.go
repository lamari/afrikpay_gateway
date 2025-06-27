package models

import (
	"fmt"
	"time"
)

// MobileMoneyConfig holds configuration for Mobile Money clients
type MobileMoneyConfig struct {
	Provider      string        // "mtn" or "orange"
	BaseURL       string
	APIKey        string
	APISecret     string
	SubscriptionKey string
	Timeout       time.Duration
	MaxRetries    int
}

// Validate validates the Mobile Money configuration
func (c *MobileMoneyConfig) Validate() error {
	if c.Provider != "mtn" && c.Provider != "orange" {
		return fmt.Errorf("provider must be 'mtn' or 'orange'")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if c.APISecret == "" {
		return fmt.Errorf("API secret is required")
	}
	if c.SubscriptionKey == "" {
		return fmt.Errorf("subscription key is required")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	return nil
}

// MTNPaymentRequest represents a payment request to MTN Mobile Money
type MTNPaymentRequest struct {
	Amount       string            `json:"amount" validate:"required"`
	Currency     string            `json:"currency" validate:"required"`
	ExternalID   string            `json:"externalId" validate:"required"`
	Payer        MTNPayer          `json:"payer" validate:"required"`
	PayerMessage string            `json:"payerMessage,omitempty"`
	PayeeNote    string            `json:"payeeNote,omitempty"`
	CallbackURL  string            `json:"callbackUrl,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// MTNPayer represents the payer information for MTN
type MTNPayer struct {
	PartyIDType string `json:"partyIdType" validate:"required,oneof=MSISDN EMAIL PARTY_CODE"`
	PartyID     string `json:"partyId" validate:"required"`
}

// MTNPaymentResponse represents the payment response from MTN
type MTNPaymentResponse struct {
	ReferenceID string `json:"referenceId"`
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
}

// MTNPaymentStatusRequest represents a payment status request
type MTNPaymentStatusRequest struct {
	ReferenceID string `json:"referenceId" validate:"required"`
}

// MTNPaymentStatusResponse represents the payment status response
type MTNPaymentStatusResponse struct {
	Amount            string            `json:"amount"`
	Currency          string            `json:"currency"`
	FinancialTransactionID string       `json:"financialTransactionId"`
	ExternalID        string            `json:"externalId"`
	Payer             MTNPayer          `json:"payer"`
	PayerMessage      string            `json:"payerMessage"`
	PayeeNote         string            `json:"payeeNote"`
	Status            string            `json:"status"`
	Reason            string            `json:"reason,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// OrangePaymentRequest represents a payment request to Orange Money
type OrangePaymentRequest struct {
	MerchantKey     string  `json:"merchant_key" validate:"required"`
	Currency        string  `json:"currency" validate:"required"`
	OrderID         string  `json:"order_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	ReturnURL       string  `json:"return_url" validate:"required,url"`
	CancelURL       string  `json:"cancel_url" validate:"required,url"`
	NotifURL        string  `json:"notif_url" validate:"required,url"`
	Lang            string  `json:"lang,omitempty"`
	Reference       string  `json:"reference,omitempty"`
	CustomerMSISDN  string  `json:"customer_msisdn,omitempty"`
}

// OrangePaymentResponse represents the payment response from Orange Money
type OrangePaymentResponse struct {
	PayToken    string `json:"pay_token"`
	PaymentURL  string `json:"payment_url"`
	NotifToken  string `json:"notif_token"`
	MerchantKey string `json:"merchant_key"`
}

// OrangePaymentStatusRequest represents a payment status request to Orange
type OrangePaymentStatusRequest struct {
	MerchantKey string `json:"merchant_key" validate:"required"`
	OrderID     string `json:"order_id" validate:"required"`
}

// OrangePaymentStatusResponse represents the payment status response from Orange
type OrangePaymentStatusResponse struct {
	Status      string  `json:"status"`
	TxnID       string  `json:"txnid"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	OrderID     string  `json:"order_id"`
	MerchantKey string  `json:"merchant_key"`
	Reference   string  `json:"reference,omitempty"`
	Message     string  `json:"message,omitempty"`
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

// PaymentStatus constants are defined in common.go
