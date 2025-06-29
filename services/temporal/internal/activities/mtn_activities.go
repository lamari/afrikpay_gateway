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
