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
		configPath = "../config/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	// Set global config for workflows/activities
	workflows.SetConfig(cfg)
	// Set activity client factories
	activities.SetConfigAwareFactories(cfg)

	c, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
	})
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "afrikpay", worker.Options{})

	// Register workflows (pass cfg if needed)
	//w.RegisterWorkflow(workflows.CreateUserWorkflow)

	// Register activities as closures that inject cfg
	//w.RegisterActivity(activities.CreateUserActivity)

	log.Println("[Temporal] Worker started on task queue: afrikpay")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("unable to start worker: %v", err)
	}
}
