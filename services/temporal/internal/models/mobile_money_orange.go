package models

// OrangePaymentRequest represents a payment request to Orange Money
type OrangePaymentRequest struct {
	MerchantKey    string  `json:"merchant_key" validate:"required"`
	Currency       string  `json:"currency" validate:"required"`
	OrderID        string  `json:"order_id" validate:"required"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	ReturnURL      string  `json:"return_url" validate:"required,url"`
	CancelURL      string  `json:"cancel_url" validate:"required,url"`
	NotifURL       string  `json:"notif_url" validate:"required,url"`
	Lang           string  `json:"lang,omitempty"`
	Reference      string  `json:"reference,omitempty"`
	CustomerMSISDN string  `json:"customer_msisdn,omitempty"`
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

