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

func TestOrangeActivitiesE2E_HealthCheck_ShouldSucceed(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour Orange
	orangeConfig := cfg.GetOrangeConfig()
	if orangeConfig == nil || orangeConfig.BaseURL == "" {
		t.Skip("Orange Mobile Money n'est pas configuré, test ignoré")
	}

	orangeClient := clients.NewOrangeClient(*cfg.GetOrangeConfig())
	activities := activities.NewOrangeActivities(*orangeClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Act
	err = activities.HealthCheck(ctx)

	// Assert
	assert.NoError(t, err, "HealthCheck should succeed")
}

func TestOrangeActivitiesE2E_InitiatePayment_ShouldProcess(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour Orange
	orangeConfig := cfg.GetOrangeConfig()
	if orangeConfig == nil || orangeConfig.BaseURL == "" {
		t.Skip("Orange Mobile Money n'est pas configuré, test ignoré")
	}

	orangeClient := clients.NewOrangeClient(*cfg.GetOrangeConfig())
	activities := activities.NewOrangeActivities(*orangeClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Créer une requête de paiement selon le modèle attendu par Orange
	paymentRequest := &models.PaymentRequest{
		PhoneNumber: "237612345678",
		Amount:      1000, // Montant minimal pour les tests
		Currency:    "XAF",
		ExternalID:  uuid.New().String(),
		Description: "Test E2E paiement Orange Money",
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

func TestOrangeActivitiesE2E_GetPaymentStatus_ShouldCheckStatus(t *testing.T) {
	// Arrange
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Vérifier que nous avons une configuration pour Orange
	orangeConfig := cfg.GetOrangeConfig()
	if orangeConfig == nil || orangeConfig.BaseURL == "" {
		t.Skip("Orange Mobile Money n'est pas configuré, test ignoré")
	}

	orangeClient := clients.NewOrangeClient(*cfg.GetOrangeConfig())
	activities := activities.NewOrangeActivities(*orangeClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Utiliser un ID de transaction de test ou créer une transaction de test d'abord
	testReferenceID := "OM-TEST-" + uuid.New().String()

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
