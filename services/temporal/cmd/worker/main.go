package main

import (
	"log"
	"os"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/workflows"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	_ = godotenv.Load() // Charge .env local si présent

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	// Set activity client factories
	activities.SetConfigAwareFactories(cfg)

	// Créer le client Temporal
	c, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Server.Address,
	})
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "afrikpay", worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.BinancePriceWorkflow)
	w.RegisterWorkflow(workflows.BinanceQuotesWorkflow)
	w.RegisterWorkflow(workflows.BinanceOrdersWorkflow)
	w.RegisterWorkflow(workflows.BinancePlaceOrderWorkflow)
	w.RegisterWorkflow(workflows.BinanceGetOrderStatusWorkflow)

	// Register activities using singleton factories
	binanceActivities := activities.GetBinanceActivitiesFromFactory()
	w.RegisterActivity(binanceActivities.GetPrice)
	w.RegisterActivity(binanceActivities.GetQuotes)
	w.RegisterActivity(binanceActivities.GetAllOrders)
	w.RegisterActivity(binanceActivities.PlaceOrder)
	w.RegisterActivity(binanceActivities.GetOrderStatus)

	log.Println("[Temporal] Worker started on task queue: afrikpay")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("unable to start worker: %v", err)
	}
}
