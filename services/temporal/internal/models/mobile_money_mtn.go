package models

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
	Amount                 string            `json:"amount"`
	Currency               string            `json:"currency"`
	FinancialTransactionID string            `json:"financialTransactionId"`
	ExternalID             string            `json:"externalId"`
	Payer                  MTNPayer          `json:"payer"`
	PayerMessage           string            `json:"payerMessage"`
	PayeeNote              string            `json:"payeeNote"`
	Status                 string            `json:"status"`
	Reason                 string            `json:"reason,omitempty"`
	Metadata               map[string]string `json:"metadata,omitempty"`
}

