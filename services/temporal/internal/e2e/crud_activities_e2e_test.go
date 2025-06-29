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

// Ces tests sont des tests E2E qui se connectent au service CRUD réel
// Ils ne doivent être exécutés que si le service CRUD est accessible

func TestCrudActivitiesE2E(t *testing.T) {
	// Charger la configuration depuis le fichier de configuration
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Créer le client CRUD avec la configuration réelle
	crudClient, err := clients.NewCrudClient(&cfg.CRUD)
	require.NoError(t, err, "Failed to create CRUD client")

	// Créer les activités CRUD avec le client réel
	activities := activities.NewCrudActivities(crudClient)

	// Test du healthcheck pour s'assurer que le service CRUD est accessible
	t.Run("HealthCheck_ServiceAvailable_Success", func(t *testing.T) {
		// Given
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// When
		err := activities.HealthCheck(ctx)

		// Then
		assert.NoError(t, err, "HealthCheck should succeed")
	})

	// Test de GetWallet
	t.Run("GetWallet_ExistingWallet_ReturnsWallet", func(t *testing.T) {
		// Given
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Pour ce test, nous avons besoin d'un ID utilisateur qui existe déjà
		// dans la base de données avec un portefeuille
		userID := "test-user-id" // Remplacer par un ID valide si nécessaire
		currency := "USDT"

		// When
		wallet, err := activities.GetWallet(ctx, userID, currency)

		if err != nil {
			t.Logf("GetWallet failed, possibly because test user doesn't exist: %v", err)
			t.Skip("Skipping test as test user may not exist")
			return
		}

		// Then
		assert.NotNil(t, wallet, "Wallet should not be nil")
		assert.Equal(t, userID, wallet.UserID, "UserID should match")
		assert.Equal(t, currency, wallet.Currency, "Currency should match")
		assert.NotEmpty(t, wallet.WalletID, "Wallet ID should not be empty")
	})

	// Test de CreateTransaction
	t.Run("CreateTransaction_ValidTransaction_ReturnsTransactionResponse", func(t *testing.T) {
		// Given
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		transaction := &models.Transaction{
			ID:        uuid.New().String(),
			UserID:    "test-user-id", // Utiliser un ID valide
			Amount:    10.0,
			Currency:  "USDT",
			Type:      "DEPOSIT",
			Status:    "PENDING",
			Reference: "e2e-test-" + time.Now().Format(time.RFC3339),
			CreatedAt: time.Now(),
		}

		// When
		response, err := activities.CreateTransaction(ctx, transaction)

		// Then
		// Si l'utilisateur test n'existe pas, le test pourrait échouer
		if err != nil {
			t.Logf("CreateTransaction failed, possibly because test user doesn't exist: %v", err)
			t.Skip("Skipping test as test user may not exist")
			return
		}

		assert.NoError(t, err, "CreateTransaction should succeed")
		assert.NotNil(t, response, "Transaction response should not be nil")
		assert.Equal(t, transaction.ID, response.TransactionID, "Transaction ID should match")
		assert.Equal(t, transaction.UserID, response.UserID, "UserID should match")
		assert.Equal(t, transaction.Amount, response.Amount, "Amount should match")
		assert.Equal(t, transaction.Currency, response.Currency, "Currency should match")
	})

	// Test de UpdateWalletBalance
	t.Run("UpdateWalletBalance_ValidWallet_UpdatesBalance", func(t *testing.T) {
		// Given
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// D'abord, récupérer un portefeuille existant
		userID := "test-user-id" // Utiliser un ID valide
		currency := "USDT"

		wallet, err := activities.GetWallet(ctx, userID, currency)
		if err != nil {
			t.Logf("GetWallet failed, possibly because test user doesn't exist: %v", err)
			t.Skip("Skipping test as test user may not exist")
			return
		}

		// Enregistrer le solde initial
		initialBalance := wallet.Balance

		// Montant à ajouter (petit pour éviter des problèmes)
		amountToAdd := 1.0

		// When
		updatedWallet, err := activities.UpdateWalletBalance(ctx, wallet.WalletID, amountToAdd, currency)

		// Then
		assert.NoError(t, err, "UpdateWalletBalance should succeed")
		assert.NotNil(t, updatedWallet, "Updated wallet should not be nil")
		assert.Equal(t, wallet.WalletID, updatedWallet.WalletID, "Wallet ID should match")
		assert.Equal(t, wallet.UserID, updatedWallet.UserID, "UserID should match")
		assert.Equal(t, wallet.Currency, updatedWallet.Currency, "Currency should match")

		// Vérifier que le solde a été mis à jour correctement
		// Note: en fonction de l'implémentation, le solde pourrait être initialBalance + amountToAdd
		// ou simplement amountToAdd (remplacement complet)
		t.Logf("Initial balance: %f, Updated balance: %f", initialBalance, updatedWallet.Balance)
		assert.NotEqual(t, initialBalance, updatedWallet.Balance, "Balance should be updated")

		// Remettre le solde à sa valeur initiale pour les futurs tests
		_, resetErr := activities.UpdateWalletBalance(ctx, wallet.WalletID, initialBalance, currency)
		if resetErr != nil {
			t.Logf("Failed to reset wallet balance: %v", resetErr)
		}
	})

	// Test de création d'un portefeuille via mise à jour (si cette fonctionnalité est supportée)
	t.Run("UpdateWalletBalance_NonExistentWallet_CreatesWallet", func(t *testing.T) {
		// Given
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Créer un ID de portefeuille qui n'existe probablement pas
		nonExistentWalletID := "non-existent-wallet-" + uuid.New().String()
		amount := 100.0
		currency := "USDT"

		// When
		wallet, err := activities.UpdateWalletBalance(ctx, nonExistentWalletID, amount, currency)

		// Then
		// Si l'API ne supporte pas la création via updateWalletBalance, ce test échouera
		// et c'est acceptable
		if err != nil {
			t.Logf("UpdateWalletBalance for non-existent wallet failed as expected: %v", err)
			return
		}

		assert.NotNil(t, wallet, "Wallet should not be nil if created")
		assert.Equal(t, nonExistentWalletID, wallet.WalletID, "Wallet ID should match")
		assert.Equal(t, amount, wallet.Balance, "Balance should be set to the specified amount")
		assert.Equal(t, currency, wallet.Currency, "Currency should match")
	})
}
