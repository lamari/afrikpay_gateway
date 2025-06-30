package workflows

import (
	"fmt"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/models"
	"go.temporal.io/sdk/workflow"
)

// MTNPaymentWorkflow orchestrates the MTN payment creation workflow
func MTNPaymentWorkflow(ctx workflow.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting MTN payment workflow", "externalId", request.ExternalID)

	// Configure activity options
	ao := defaultActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Create user and reference ID
	var referenceID string
	err := workflow.ExecuteActivity(ctx, activities.GetMTNActivitiesFromFactory().CreateUser, request.CallbackURL).Get(ctx, &referenceID)
	if err != nil {
		logger.Error("CreateUser failed", "error", err)
		return nil, fmt.Errorf("CreateUser error: %w", err)
	}

	// Step 2: Create API key
	var apiKey string
	err = workflow.ExecuteActivity(ctx, activities.GetMTNActivitiesFromFactory().CreateApiKey, referenceID).Get(ctx, &apiKey)
	if err != nil {
		logger.Error("CreateApiKey failed", "error", err)
		return nil, fmt.Errorf("CreateApiKey error: %w", err)
	}

	// Step 3: Get access token
	var accessToken string
	err = workflow.ExecuteActivity(ctx, activities.GetMTNActivitiesFromFactory().GetAccessToken, referenceID, apiKey).Get(ctx, &accessToken)
	if err != nil {
		logger.Error("GetAccessToken failed", "error", err)
		return nil, fmt.Errorf("GetAccessToken error: %w", err)
	}

	// Step 4: Create payment request
	mtnReq := models.MTNPaymentRequest{
		Amount:      fmt.Sprintf("%.2f", request.Amount),
		Currency:    request.Currency,
		ExternalID:  request.ExternalID,
		Payer:       models.MTNPayer{PartyIDType: "MSISDN", PartyID: request.PhoneNumber},
		PayerMessage: request.Description,
		PayeeNote:   request.Description,
		CallbackURL: request.CallbackURL,
		Metadata:    request.Metadata,
	}
	var mtnResp *models.MTNPaymentResponse
	err = workflow.ExecuteActivity(ctx, activities.GetMTNActivitiesFromFactory().CreatePaymentRequest, referenceID, accessToken, &mtnReq).Get(ctx, &mtnResp)
	if err != nil {
		logger.Error("CreatePaymentRequest failed", "error", err)
		return nil, fmt.Errorf("CreatePaymentRequest error: %w", err)
	}

	// Map MTN response to generic PaymentResponse
	result := &models.PaymentResponse{
		ReferenceID: referenceID,
		Status:      models.PaymentStatus(mtnResp.Status),
		PaymentURL:  "",
		Message:     mtnResp.Reason,
	}

	logger.Info("MTN payment workflow completed", "referenceId", referenceID, "status", result.Status)
	return result, nil
}
