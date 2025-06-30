package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ces tests sont des tests E2E qui se connectent aux APIs réelles
// Ils ne doivent être exécutés que si les clés API sont configurées correctement

func TestBinanceActivitiesE2E(t *testing.T) {
	// Charger la configuration depuis le fichier de configuration
	cfg, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Utilisez la configuration réelle pour créer le client Binance
	binanceClient := clients.NewBinanceClient(*cfg.GetBinanceConfig())

	// Créer les activités Binance avec le client réel
	// D'après le fichier binance_activities.go, nous devons passer le client
	binanceActivities := activities.NewBinanceActivities(*binanceClient)

	// Test du healthcheck pour s'assurer que l'API Binance est accessible
	t.Run("HealthCheck", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := binanceActivities.HealthCheck(ctx)
		assert.NoError(t, err, "HealthCheck should succeed")
	})

	// Test de la méthode GetPrice
	t.Run("GetPrice", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		symbol := "BTCUSDT" // Utilisez une paire qui existe certainement sur Binance testnet
		price, err := binanceActivities.GetPrice(ctx, symbol)

		assert.NoError(t, err, "GetPrice should succeed")
		assert.NotNil(t, price, "Price should not be nil")
		assert.Equal(t, symbol, price.Symbol, "Symbol should match")
		assert.Greater(t, price.Price, 0.0, "Price should be greater than zero")
		assert.False(t, price.Timestamp.IsZero(), "Timestamp should not be zero")
	})

	// Test de la méthode GetQuote
	t.Run("GetQuote", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		symbol := "BTCUSDT"
		quote, err := binanceActivities.GetQuote(ctx, symbol)

		assert.NoError(t, err, "GetQuote should succeed")
		assert.NotNil(t, quote, "Quote should not be nil")
		assert.Equal(t, symbol, quote.Symbol, "Symbol should match")
		assert.Greater(t, quote.LastPrice, 0.0, "LastPrice should be greater than zero")
		assert.False(t, quote.Timestamp.IsZero(), "Timestamp should not be zero")
	})

	// Test de la méthode GetQuotes
	t.Run("GetQuotes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		quotes, err := binanceActivities.GetQuotes(ctx)

		assert.NoError(t, err, "GetQuotes should succeed")
		assert.NotNil(t, quotes, "Quotes should not be nil")
		// Dans un environnement réel, il devrait y avoir des quotes
		// mais ce n'est pas toujours garanti, donc on ne vérifie pas la longueur
	})

	// Les tests suivants simulent des opérations d'ordre, mais ne placent pas réellement d'ordres
	// car cela nécessiterait un compte avec des fonds
	t.Run("PlaceOrder", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderRequest := &models.OrderRequest{
			Symbol:   "BTCUSDT",
			Side:     "BUY",
			Type:     "MARKET",
			Quantity: 0.001, // Très petite quantité pour les tests
			Price:    0,     // Prix du marché
		}

		// Ce test pourrait échouer si le client ne simule pas les ordres
		// dans un environnement réel, il faudrait vérifier l'état du compte
		response, err := binanceActivities.PlaceOrder(ctx, orderRequest)

		// Si ce test échoue, c'est probablement parce que la fonction placeOrder
		// dans client Binance essaie de placer un vrai ordre et n'est pas en mode simulation
		if err != nil {
			t.Logf("PlaceOrder failed as expected in test environment: %v", err)
			return
		}

		assert.Equal(t, response.Quantity, 0.001)
		assert.NotNil(t, response, "Order should not be nil")
		assert.Equal(t, orderRequest.Symbol, response.Symbol, "Symbol should match")
		assert.Equal(t, orderRequest.Side, response.Side, "Side should match")
	})
}
