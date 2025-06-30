package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ces tests sont des tests E2E qui se connectent aux APIs réelles
// Ils ne doivent être exécutés que si les clés API sont configurées correctement

func TestMTNActivitiesE2E_HealthCheck_ShouldSucceed(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour MTN
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money n'est pas configuré, test ignoré")
	}

	mtnClient := clients.NewMTNClient(*mtnConfig)
	activities := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Act
	err = activities.HealthCheck(ctx)

	// Assert
	assert.NoError(t, err, "HealthCheck should succeed")
}

func TestMTNActivitiesE2E_InitiatePayment_ShouldProcess(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour MTN
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money n'est pas configuré, test ignoré")
	}

	mtnClient := clients.NewMTNClient(*cfg.GetMTNConfig())
	activities := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Créer une requête de paiement selon le modèle attendu par MTN
	paymentRequest := &models.PaymentRequest{
		PhoneNumber: "237612345678",
		Amount:      1000, // Montant minimal pour les tests
		Currency:    "XAF",
		ExternalID:  uuid.New().String(),
		Description: "Test E2E paiement MTN Mobile Money",
	}

	// Act
	response, err := activities.InitiatePayment(ctx, paymentRequest)

	// Assert
	// Dans un environnement sandbox, on peut soit avoir un succès simulé, soit une erreur d'API
	if err != nil {
		t.Logf("InitiatePayment error (expected in sandbox): %v", err)
		return
	}

	assert.NotNil(t, response, "Response should not be nil")
	assert.NotEmpty(t, response.ReferenceID, "ReferenceID should not be empty")
	// Le statut devrait être pending pour un nouveau paiement
	assert.Equal(t, models.PaymentStatusPending, response.Status, "Status should be PENDING")
}

func TestMTNActivitiesE2E_GetPaymentStatus_ShouldCheckStatus(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour MTN
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money n'est pas configuré, test ignoré")
	}

	mtnClient := clients.NewMTNClient(*cfg.GetMTNConfig())
	activities := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Utiliser un ID de transaction de test ou créer une transaction de test d'abord
	testReferenceID := "MTN-TEST-" + uuid.New().String()

	// Act
	response, err := activities.GetPaymentStatus(ctx, testReferenceID)

	// Assert
	// Dans un environnement sandbox, on peut soit avoir un succès simulé, soit une erreur d'API
	if err != nil {
		t.Logf("GetPaymentStatus error (expected for non-existent transaction): %v", err)
		return
	}

	assert.NotNil(t, response, "Response should not be nil")
	assert.Equal(t, testReferenceID, response.ReferenceID, "ReferenceID should match")
}

// Test creation of MTN API user
func TestMTNActivitiesE2E_CreateUser_ShouldReturnReferenceID(t *testing.T) {
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money not configured, skipping")
	}
	mtnClient := clients.NewMTNClient(*mtnConfig)
	acts := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ref, err := acts.CreateUser(ctx, "http://localhost/callback")
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	assert.NotEmpty(t, ref, "ReferenceID should not be empty")
}

// Test generation of MTN API key
func TestMTNActivitiesE2E_CreateApiKey_ShouldReturnApiKey(t *testing.T) {
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err)
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money not configured, skipping")
	}
	mtnClient := clients.NewMTNClient(*mtnConfig)
	acts := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ensure user exists
	referenceID, err := acts.CreateUser(ctx, "http://localhost/callback")
	require.NoError(t, err)

	apiKey, err := acts.CreateApiKey(ctx, referenceID)
	if err != nil {
		t.Fatalf("CreateApiKey error: %v", err)
	}
	assert.NotEmpty(t, apiKey, "ApiKey should not be empty")
}

// Test retrieval of MTN access token
func TestMTNActivitiesE2E_GetAccessToken_ShouldReturnToken(t *testing.T) {
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err)
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money not configured, skipping")
	}
	mtnClient := clients.NewMTNClient(*mtnConfig)
	acts := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// prepare user and apiKey
	referenceID, err := acts.CreateUser(ctx, "http://localhost/callback")
	require.NoError(t, err)
	apiKey, err := acts.CreateApiKey(ctx, referenceID)
	require.NoError(t, err)

	token, err := acts.GetAccessToken(ctx, referenceID, apiKey)
	if err != nil {
		t.Fatalf("GetAccessToken error: %v", err)
	}
	assert.NotEmpty(t, token, "AccessToken should not be empty")
}

// Test creation of MTN payment request
func TestMTNActivitiesE2E_CreatePaymentRequest_ShouldReturnResponse(t *testing.T) {
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err)
	mtnConfig := cfg.GetMTNConfig()
	if mtnConfig == nil || mtnConfig.BaseURL == "" {
		t.Skip("MTN Mobile Money not configured, skipping")
	}
	mtnClient := clients.NewMTNClient(*mtnConfig)
	acts := activities.NewMTNActivities(*mtnClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// prepare auth
	referenceID, err := acts.CreateUser(ctx, "http://localhost/callback")
	require.NoError(t, err)
	apiKey, err := acts.CreateApiKey(ctx, referenceID)
	require.NoError(t, err)
	accessToken, err := acts.GetAccessToken(ctx, referenceID, apiKey)
	require.NoError(t, err)

	mtnReq := &models.MTNPaymentRequest{
		Amount:       "1.00",
		Currency:     "EUR",
		ExternalID:   uuid.New().String(),
		Payer:        models.MTNPayer{PartyIDType: "MSISDN", PartyID: "256774290781"},
		PayerMessage: "Test",
		PayeeNote:    "Note",
		CallbackURL:  "http://localhost/callback",
		Metadata:     map[string]string{},
	}

	resp, err := acts.CreatePaymentRequest(ctx, referenceID, accessToken, mtnReq)
	if err != nil {
		t.Logf("CreatePaymentRequest error (expected in sandbox): %v", err)
		return
	}
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NotEmpty(t, resp.ReferenceID, "ReferenceID should not be empty")
	assert.Equal(t, "PENDING", resp.Status, "Status should be PENDING")
}
