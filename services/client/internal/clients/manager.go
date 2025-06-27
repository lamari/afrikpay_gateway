package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/afrikpay/gateway/services/client/internal/models"
)

// ClientManager manages all external API clients
type ClientManager struct {
	binanceClient     ExchangeClient
	bitgetClient      ExchangeClient
	mtnClient         MobileMoneyClient
	orangeClient      MobileMoneyClient
	mu                sync.RWMutex
}

// ExchangeClient defines the interface for cryptocurrency exchange clients
type ExchangeClient interface {
	GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error)
	PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error)
	GetOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error)
	GetQuotes(ctx context.Context) (*models.QuotesResponse, error)
	GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error)
	GetResilienceStats() *models.ResilienceStats
	ResetResilienceStats()
	Close() error
}

// MobileMoneyClient defines the interface for mobile money clients
type MobileMoneyClient interface {
	models.ClientInterface
	
	// Initiate payment
	InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error)
	
	// Get payment status
	GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error)
}

// NewClientManager creates a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{}
}

// SetBinanceClient sets the Binance client
func (cm *ClientManager) SetBinanceClient(client ExchangeClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.binanceClient = client
}

// SetBitgetClient sets the Bitget client
func (cm *ClientManager) SetBitgetClient(client ExchangeClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.bitgetClient = client
}

// SetMTNClient sets the MTN Mobile Money client
func (cm *ClientManager) SetMTNClient(client MobileMoneyClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.mtnClient = client
}

// SetOrangeClient sets the Orange Mobile Money client
func (cm *ClientManager) SetOrangeClient(client MobileMoneyClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.orangeClient = client
}

// GetBinanceClient returns the Binance client
func (cm *ClientManager) GetBinanceClient() (ExchangeClient, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if cm.binanceClient == nil {
		return nil, fmt.Errorf("Binance client not initialized")
	}
	return cm.binanceClient, nil
}

// GetBitgetClient returns the Bitget client
func (cm *ClientManager) GetBitgetClient() (ExchangeClient, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if cm.bitgetClient == nil {
		return nil, fmt.Errorf("Bitget client not initialized")
	}
	return cm.bitgetClient, nil
}

// GetMTNClient returns the MTN Mobile Money client
func (cm *ClientManager) GetMTNClient() (MobileMoneyClient, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if cm.mtnClient == nil {
		return nil, fmt.Errorf("MTN client not initialized")
	}
	return cm.mtnClient, nil
}

// GetOrangeClient returns the Orange Mobile Money client
func (cm *ClientManager) GetOrangeClient() (MobileMoneyClient, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if cm.orangeClient == nil {
		return nil, fmt.Errorf("Orange client not initialized")
	}
	return cm.orangeClient, nil
}

// GetExchangeClient returns an exchange client by name
func (cm *ClientManager) GetExchangeClient(exchange string) (ExchangeClient, error) {
	switch exchange {
	case "binance":
		return cm.GetBinanceClient()
	case "bitget":
		return cm.GetBitgetClient()
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", exchange)
	}
}

// GetMobileMoneyClient returns a mobile money client by provider
func (cm *ClientManager) GetMobileMoneyClient(provider string) (MobileMoneyClient, error) {
	switch provider {
	case "mtn":
		return cm.GetMTNClient()
	case "orange":
		return cm.GetOrangeClient()
	default:
		return nil, fmt.Errorf("unsupported mobile money provider: %s", provider)
	}
}

// GetAllResilienceStats returns resilience statistics for all clients
func (cm *ClientManager) GetAllResilienceStats() map[string]*models.ResilienceStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]*models.ResilienceStats)

	if cm.binanceClient != nil {
		stats["binance"] = cm.binanceClient.GetResilienceStats()
	}

	if cm.bitgetClient != nil {
		stats["bitget"] = cm.bitgetClient.GetResilienceStats()
	}

	if cm.mtnClient != nil {
		stats["mtn"] = cm.mtnClient.GetResilienceStats()
	}

	if cm.orangeClient != nil {
		stats["orange"] = cm.orangeClient.GetResilienceStats()
	}

	return stats
}

// ResetAllResilienceStats resets resilience statistics for all clients
func (cm *ClientManager) ResetAllResilienceStats() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.binanceClient != nil {
		cm.binanceClient.ResetResilienceStats()
	}

	if cm.bitgetClient != nil {
		cm.bitgetClient.ResetResilienceStats()
	}

	if cm.mtnClient != nil {
		cm.mtnClient.ResetResilienceStats()
	}

	if cm.orangeClient != nil {
		cm.orangeClient.ResetResilienceStats()
	}
}

// HealthCheck performs health checks on all clients
func (cm *ClientManager) HealthCheck(ctx context.Context) map[string]bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	health := make(map[string]bool)

	// Check Binance client
	if cm.binanceClient != nil {
		_, err := cm.binanceClient.GetQuotes(ctx)
		health["binance"] = err == nil
	} else {
		health["binance"] = false
	}

	// Check Bitget client
	if cm.bitgetClient != nil {
		_, err := cm.bitgetClient.GetQuotes(ctx)
		health["bitget"] = err == nil
	} else {
		health["bitget"] = false
	}

	// Check MTN client - we can't easily test without making actual API calls
	// so we just check if the client is initialized
	health["mtn"] = cm.mtnClient != nil

	// Check Orange client - same as MTN
	health["orange"] = cm.orangeClient != nil

	return health
}

// Close closes all clients and releases resources
func (cm *ClientManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var errors []error

	if cm.binanceClient != nil {
		if err := cm.binanceClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Binance client: %w", err))
		}
	}

	if cm.bitgetClient != nil {
		if err := cm.bitgetClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Bitget client: %w", err))
		}
	}

	if cm.mtnClient != nil {
		if err := cm.mtnClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close MTN client: %w", err))
		}
	}

	if cm.orangeClient != nil {
		if err := cm.orangeClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Orange client: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing clients: %v", errors)
	}

	return nil
}

// GetSupportedExchanges returns a list of supported exchanges
func (cm *ClientManager) GetSupportedExchanges() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var exchanges []string
	if cm.binanceClient != nil {
		exchanges = append(exchanges, "binance")
	}
	if cm.bitgetClient != nil {
		exchanges = append(exchanges, "bitget")
	}
	return exchanges
}

// GetSupportedMobileMoneyProviders returns a list of supported mobile money providers
func (cm *ClientManager) GetSupportedMobileMoneyProviders() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var providers []string
	if cm.mtnClient != nil {
		providers = append(providers, "mtn")
	}
	if cm.orangeClient != nil {
		providers = append(providers, "orange")
	}
	return providers
}
