package main

import (
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
	cfg := config.MustLoadConfig(configPath)

	e := echo.New()

	// Initialisation du client Temporal (production)
	c, err := goTemporalClient.Dial(goTemporalClient.Options{
		HostPort: cfg.Temporal.Address,
	})
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	handler.SetTemporalClient(c)

	handler.RegisterRoutes(e, cfg)

	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
