package workflows

import (
	"fmt"
	"strings"
	"time"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/models"
	"go.temporal.io/sdk/temporal"
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
	// Temporary fix: use hardcoded callback URL for testing
	callbackURL := "http://localhost/callback"
	if request.CallbackURL != "" {
		callbackURL = request.CallbackURL
	}
	err := workflow.ExecuteActivity(ctx, activities.GetMTNActivitiesFromFactory().CreateUser, callbackURL).Get(ctx, &referenceID)
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
	// Ensure ExternalID is not empty (required by MTN API)
	externalID := request.ExternalID
	if externalID == "" {
		externalID = "mtn-" + workflow.GetInfo(ctx).WorkflowExecution.ID[:8] + "-" + fmt.Sprintf("%d", workflow.Now(ctx).Unix())
	}

	// Ensure required fields have values
	payerMessage := request.Description
	if payerMessage == "" {
		payerMessage = "Payment request"
	}
	payeeNote := request.Description
	if payeeNote == "" {
		payeeNote = "Payment from Afrikpay"
	}
	callbackURL = request.CallbackURL
	if callbackURL == "" {
		callbackURL = "http://localhost/callback"
	}

	mtnReq := models.MTNPaymentRequest{
		Amount:       fmt.Sprintf("%.2f", request.Amount),
		Currency:     request.Currency,
		ExternalID:   externalID,
		Payer:        models.MTNPayer{PartyIDType: "MSISDN", PartyID: "256774290781"},
		PayerMessage: payerMessage,
		PayeeNote:    payeeNote,
		CallbackURL:  callbackURL,
		Metadata:     request.Metadata,
	}

	// Configurer des options spécifiques pour CreatePaymentRequest - sans retry
	createPaymentRequestOptions := workflow.ActivityOptions{
		StartToCloseTimeout:    10 * time.Second,
		ScheduleToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1, // Exécuter une seule fois, pas de retry
		},
	}

	// Créer un nouveau contexte avec les options spécifiques pour cette activité uniquement
	createPaymentCtx := workflow.WithActivityOptions(ctx, createPaymentRequestOptions)

	// Exécuter l'activité une seule fois
	var mtnResp *models.MTNPaymentResponse
	err = workflow.ExecuteActivity(createPaymentCtx, activities.GetMTNActivitiesFromFactory().CreatePaymentRequest, referenceID, accessToken, &mtnReq).Get(ctx, &mtnResp)
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

	// Step 5: Enregistrer la transaction dans le service CRUD
	transactionType := models.TransactionTypeDeposit
	transactionStatus := models.TransactionStatusPending
	if mtnResp.Status == "SUCCESSFUL" {
		transactionStatus = models.TransactionStatusCompleted
	}

	transaction := &models.Transaction{
		ID:       generateID(),
		WalletID: getWalletID(),
		UserID:   generateID(),
		Type:     strings.ToLower(string(transactionType)),
		Status:   strings.ToLower(string(transactionStatus)),
		Amount:   request.Amount,
		Currency: request.Currency,
	}

	var transactionResp *models.TransactionResponse
	err = workflow.ExecuteActivity(ctx, activities.GetCrudActivitiesFromFactory().CreateTransaction, transaction).Get(ctx, &transactionResp)
	if err != nil {
		logger.Error("CreateTransaction failed", "error", err)
		// On continue même en cas d'erreur, car le paiement MTN a déjà été créé
		// On log simplement l'erreur
		logger.Info("Failed to record transaction in CRUD service, but MTN payment was created", "error", err)
	} else {
		logger.Info("Transaction recorded successfully", "transactionId", transactionResp.TransactionID)
		// Ajouter une note sur l'ID de transaction dans le message
		if result.Message == "" {
			result.Message = fmt.Sprintf("Transaction ID: %s", transactionResp.TransactionID)
		} else {
			result.Message += fmt.Sprintf(", Transaction ID: %s", transactionResp.TransactionID)
		}
	}

	logger.Info("MTN payment workflow completed", "referenceId", referenceID, "status", result.Status)
	return result, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UTC().Unix())
}

func getWalletID() string {
	// TODO: Get wallet ID from database
	return "12345"
}
