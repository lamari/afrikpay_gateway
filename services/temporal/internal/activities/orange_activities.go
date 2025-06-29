package activities

import (
	"context"

	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/models"
)

// OrangeActivities wraps Orange Money client functionality for Temporal activities
type OrangeActivities struct {
	client clients.OrangeClient
}

// NewOrangeActivities creates a new OrangeActivities instance
func NewOrangeActivities(client clients.OrangeClient) *OrangeActivities {
	return &OrangeActivities{
		client: client,
	}
}

// InitiatePayment initiates a payment with Orange Money
func (a *OrangeActivities) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	return a.client.InitiatePayment(ctx, request)
}

// GetPaymentStatus gets the status of a payment
func (a *OrangeActivities) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	return a.client.GetPaymentStatus(ctx, referenceID)
}

// HealthCheck performs a health check against the Orange Money API
func (a *OrangeActivities) HealthCheck(ctx context.Context) error {
	return a.client.HealthCheck(ctx)
}
