package workflows

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/models"
)

func TestMTNPaymentWorkflow_Success(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	// Prepare input request
	req := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "EUR",
		ExternalID:  "ext123",
		PhoneNumber: "256774290781",
		Description: "Test Payment",
		CallbackURL: "http://localhost/callback",
		Metadata:    map[string]string{"env": "test"},
	}

	// Expected MTN request
	mtnReq := models.MTNPaymentRequest{
		Amount:       fmt.Sprintf("%.2f", req.Amount),
		Currency:     req.Currency,
		ExternalID:   req.ExternalID,
		Payer:        models.MTNPayer{PartyIDType: "MSISDN", PartyID: req.PhoneNumber},
		PayerMessage: req.Description,
		PayeeNote:    req.Description,
		CallbackURL:  req.CallbackURL,
		Metadata:     req.Metadata,
	}

	// Mock activities
	expectedRef := "test-ref-id"
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreateUser, mock.Anything, req.CallbackURL).Return(expectedRef, nil)
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreateApiKey, mock.Anything, expectedRef).Return("api-key", nil)
	env.OnActivity(activities.GetMTNActivitiesFromFactory().GetAccessToken, mock.Anything, expectedRef, "api-key").Return("access-token", nil)
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreatePaymentRequest, mock.Anything, expectedRef, "access-token", &mtnReq).Return(&models.MTNPaymentResponse{
		ReferenceID: expectedRef,
		Status:      "PENDING",
		Reason:      "ok",
	}, nil)

	// Execute workflow
	env.ExecuteWorkflow(MTNPaymentWorkflow, req)

	// Assertions
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())

	var result models.PaymentResponse
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
	assert.Equal(t, expectedRef, result.ReferenceID)
	assert.Equal(t, models.PaymentStatus("PENDING"), result.Status)
	assert.Equal(t, "ok", result.Message)
}

func TestMTNPaymentWorkflow_CreateUserError(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	// Simulate CreateUser failure
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreateUser, mock.Anything, "http://localhost/callback").Return("", fmt.Errorf("create user failed"))

	// Execute workflow
	env.ExecuteWorkflow(MTNPaymentWorkflow, &models.PaymentRequest{CallbackURL: "http://localhost/callback"})

	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())
}

func TestMTNPaymentWorkflow_PaymentRequestError(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	// Mock intermediate successes
	expectedRef := "test-ref-id"
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreateUser, mock.Anything, "http://localhost/callback").Return(expectedRef, nil)
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreateApiKey, mock.Anything, expectedRef).Return("api-key", nil)
	env.OnActivity(activities.GetMTNActivitiesFromFactory().GetAccessToken, mock.Anything, expectedRef, "api-key").Return("access-token", nil)

	// Simulate CreatePaymentRequest failure
	env.OnActivity(activities.GetMTNActivitiesFromFactory().CreatePaymentRequest, mock.Anything, expectedRef, "access-token", mock.Anything).Return(nil, fmt.Errorf("payment failed"))

	// Execute workflow
	env.ExecuteWorkflow(MTNPaymentWorkflow, &models.PaymentRequest{CallbackURL: "http://localhost/callback"})

	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())
}
