package activities

import (
	"context"

	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/models"
)

// MTNActivities wraps MTN Mobile Money client functionality for Temporal activities
type MTNActivities struct {
	client clients.MTNClient
}

// NewMTNActivities creates a new MTNActivities instance
func NewMTNActivities(client clients.MTNClient) *MTNActivities {
	return &MTNActivities{
		client: client,
	}
}

// InitiatePayment initiates a payment with MTN Mobile Money
func (a *MTNActivities) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	return a.client.InitiatePayment(ctx, request)
}

// GetPaymentStatus gets the status of a payment
func (a *MTNActivities) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	return a.client.GetPaymentStatus(ctx, referenceID)
}

// HealthCheck performs a health check against the MTN Mobile Money API
func (a *MTNActivities) HealthCheck(ctx context.Context) error {
	return a.client.HealthCheck(ctx)
}

// CreateUser creates a new API user in MTN MoMo
func (a *MTNActivities) CreateUser(ctx context.Context, providerCallbackHost string) (string, error) {
	return a.client.CreateUser(ctx, providerCallbackHost)
}

// CreateApiKey generates an API key for the MTN API user
func (a *MTNActivities) CreateApiKey(ctx context.Context, referenceID string) (string, error) {
	return a.client.CreateApiKey(ctx, referenceID)
}

// GetAccessToken retrieves an access token using the MTN API key
func (a *MTNActivities) GetAccessToken(ctx context.Context, referenceID, apiKey string) (string, error) {
	return a.client.GetAccessToken(ctx, referenceID, apiKey)
}

// CreatePaymentRequest sends a payment request and returns the MTN response
func (a *MTNActivities) CreatePaymentRequest(ctx context.Context, referenceID, accessToken string, request *models.MTNPaymentRequest) (*models.MTNPaymentResponse, error) {
	return a.client.CreatePaymentRequest(ctx, referenceID, accessToken, request)
}
