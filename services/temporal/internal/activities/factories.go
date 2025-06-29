package activities

import (
	"log"

	"github.com/afrikpay/gateway/internal/clients"
	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/models"
)

// Global client factories - these will be initialized once with configuration
var (
	binanceClientFactory func() *clients.BinanceClient
	bitgetClientFactory  func() *clients.BitgetClient
	mtnClientFactory     func() *clients.MTNClient
	orangeClientFactory  func() *clients.OrangeClient
	crudClientFactory    func() models.CrudClient
)

// Singleton instances for activities
var (
	binanceActivitiesInstance *BinanceActivities
	bitgetActivitiesInstance  *BitgetActivities
	mtnActivitiesInstance     *MTNActivities
	orangeActivitiesInstance  *OrangeActivities
	crudActivitiesInstance    *CrudActivities
)

// SetConfigAwareFactories initializes client factories with configuration
// This should be called once at application startup
func SetConfigAwareFactories(cfg *config.Config) {
	// Initialize Binance client factory
	binanceClientFactory = func() *clients.BinanceClient {
		return clients.NewBinanceClient(*cfg.GetBinanceConfig())
	}

	// Initialize Bitget client factory
	bitgetClientFactory = func() *clients.BitgetClient {
		client, err := clients.NewBitgetClient(cfg.GetBitgetConfig())
		if err != nil {
			log.Fatalf("Failed to create Bitget client: %v", err)
		}
		return client
	}

	// Initialize MTN client factory
	mtnClientFactory = func() *clients.MTNClient {
		return clients.NewMTNClient(*cfg.GetMTNConfig())
	}

	// Initialize Orange client factory
	orangeClientFactory = func() *clients.OrangeClient {
		return clients.NewOrangeClient(*cfg.GetOrangeConfig())
	}

	// Initialize CRUD client factory
	crudClientFactory = func() models.CrudClient {
		client, err := clients.NewCrudClient(&cfg.CRUD)
		if err != nil {
			log.Fatalf("Failed to create CRUD client: %v", err)
		}
		return client
	}
}

// Singleton factory functions to get activity instances
func GetBinanceActivitiesFromFactory() *BinanceActivities {
	if binanceActivitiesInstance == nil {
		if binanceClientFactory == nil {
			panic("Binance client factory not initialized. Call SetConfigAwareFactories first.")
		}
		binanceActivitiesInstance = NewBinanceActivities(*binanceClientFactory())
	}
	return binanceActivitiesInstance
}

func GetBitgetActivitiesFromFactory() *BitgetActivities {
	if bitgetActivitiesInstance == nil {
		if bitgetClientFactory == nil {
			panic("Bitget client factory not initialized. Call SetConfigAwareFactories first.")
		}
		bitgetActivitiesInstance = NewBitgetActivities(*bitgetClientFactory())
	}
	return bitgetActivitiesInstance
}

func GetMTNActivitiesFromFactory() *MTNActivities {
	if mtnActivitiesInstance == nil {
		if mtnClientFactory == nil {
			panic("MTN client factory not initialized. Call SetConfigAwareFactories first.")
		}
		mtnActivitiesInstance = NewMTNActivities(*mtnClientFactory())
	}
	return mtnActivitiesInstance
}

func GetOrangeActivitiesFromFactory() *OrangeActivities {
	if orangeActivitiesInstance == nil {
		if orangeClientFactory == nil {
			panic("Orange client factory not initialized. Call SetConfigAwareFactories first.")
		}
		orangeActivitiesInstance = NewOrangeActivities(*orangeClientFactory())
	}
	return orangeActivitiesInstance
}

func GetCrudActivitiesFromFactory() *CrudActivities {
	if crudActivitiesInstance == nil {
		if crudClientFactory == nil {
			panic("CRUD client factory not initialized. Call SetConfigAwareFactories first.")
		}
		crudActivitiesInstance = NewCrudActivities(crudClientFactory())
	}
	return crudActivitiesInstance
}
