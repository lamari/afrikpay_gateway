package main

import (
	"log"
	"os"

	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/workflows"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/activity"
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

	// Create worker with the task queue "afrikpay"
	w := worker.New(c, "afrikpay", worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.BinancePriceWorkflow)
	w.RegisterWorkflow(workflows.BinanceQuotesWorkflow)
	w.RegisterWorkflow(workflows.BinanceOrdersWorkflow)
	w.RegisterWorkflow(workflows.BinancePlaceOrderWorkflow)
	w.RegisterWorkflow(workflows.BinanceGetOrderStatusWorkflow)
	w.RegisterWorkflow(workflows.MTNPaymentWorkflow)

	// Register Binance activities using singleton factory
	binanceActivities := activities.GetBinanceActivitiesFromFactory()
	w.RegisterActivityWithOptions(binanceActivities.GetPrice, activity.RegisterOptions{Name: "GetPrice"})
	w.RegisterActivityWithOptions(binanceActivities.GetQuotes, activity.RegisterOptions{Name: "GetQuotes"})
	w.RegisterActivityWithOptions(binanceActivities.GetAllOrders, activity.RegisterOptions{Name: "GetAllOrders"})
	w.RegisterActivityWithOptions(binanceActivities.PlaceOrder, activity.RegisterOptions{Name: "PlaceOrder"})
	w.RegisterActivityWithOptions(binanceActivities.GetOrderStatus, activity.RegisterOptions{Name: "GetOrderStatus"})

	// Register MTN activities using singleton factory
	mtnActivities := activities.GetMTNActivitiesFromFactory()
	w.RegisterActivityWithOptions(mtnActivities.InitiatePayment, activity.RegisterOptions{Name: "InitiatePayment"})
	w.RegisterActivityWithOptions(mtnActivities.GetPaymentStatus, activity.RegisterOptions{Name: "GetPaymentStatus"})
	w.RegisterActivityWithOptions(mtnActivities.HealthCheck, activity.RegisterOptions{Name: "HealthCheck"})
	w.RegisterActivityWithOptions(mtnActivities.CreateUser, activity.RegisterOptions{Name: "CreateUser"})
	w.RegisterActivityWithOptions(mtnActivities.CreateApiKey, activity.RegisterOptions{Name: "CreateApiKey"})
	w.RegisterActivityWithOptions(mtnActivities.GetAccessToken, activity.RegisterOptions{Name: "GetAccessToken"})
	w.RegisterActivityWithOptions(mtnActivities.CreatePaymentRequest, activity.RegisterOptions{Name: "CreatePaymentRequest"})

	// Register CRUD activities if available
	crudActivities := activities.GetCrudActivitiesFromFactory()
	if crudActivities != nil {
		w.RegisterActivityWithOptions(crudActivities.CreateTransaction, activity.RegisterOptions{Name: "CreateTransaction"})
		w.RegisterActivityWithOptions(crudActivities.UpdateWalletBalance, activity.RegisterOptions{Name: "UpdateWalletBalance"})
		w.RegisterActivityWithOptions(crudActivities.GetWallet, activity.RegisterOptions{Name: "GetWallet"})
		w.RegisterActivityWithOptions(crudActivities.HealthCheck, activity.RegisterOptions{Name: "CrudHealthCheck"})
	}

	log.Println("[Temporal] Worker started on task queue: afrikpay")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("unable to start worker: %v", err)
	}
}
