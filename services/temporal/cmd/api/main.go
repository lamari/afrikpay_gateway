package main

import (
	"fmt"
	"log"
	"os"

	"github.com/afrikpay/gateway/internal/config"
	"github.com/afrikpay/gateway/internal/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	goTemporalClient "go.temporal.io/sdk/client"
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

	e := echo.New()

	// Initialisation du client Temporal (production)
	c, err := goTemporalClient.Dial(goTemporalClient.Options{
		HostPort: cfg.Temporal.Server.Address,
	})
	if err != nil {
		log.Printf("Warning: unable to create Temporal client: %v", err)
		log.Printf("Continuing without Temporal client...")
	} else {
		defer c.Close()
		handler.SetTemporalClient(c)
	}

	handler.RegisterRoutes(e, cfg)

	port := fmt.Sprintf(":%d", cfg.Temporal.API.Port)
	log.Printf("Starting server on port %s (PortAPI value: %d)", port, cfg.Temporal.API.Port)
	e.Logger.Fatal(e.Start(port))
}
