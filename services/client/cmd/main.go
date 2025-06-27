package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/clients"
	"github.com/afrikpay/gateway/services/client/internal/config"
	"github.com/afrikpay/gateway/services/client/internal/handlers"
	"github.com/afrikpay/gateway/services/client/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../../config/config.yml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize clients
	clientManager, err := initializeClients(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}
	defer clientManager.Close()

	// Initialize handlers
	handler := handlers.NewClientHandler(clientManager)

	// Setup router
	router := setupRouter(handler)

	// Setup server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting client service on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// initializeClients initializes all external API clients
func initializeClients(cfg *config.Config) (*clients.ClientManager, error) {
	manager := clients.NewClientManager()

	// Initialize Binance client
	binanceClient := clients.NewBinanceClient(*cfg.GetBinanceConfig())
	manager.SetBinanceClient(binanceClient)

	// Initialize Bitget client
	bitgetClient, err := clients.NewBitgetClient(cfg.GetBitgetConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create Bitget client: %w", err)
	}
	manager.SetBitgetClient(bitgetClient)

	// Initialize MTN Mobile Money client
	mtnClient, err := clients.NewMobileMoneyClient(*cfg.GetMTNConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create MTN client: %w", err)
	}
	manager.SetMTNClient(mtnClient)

	// Initialize Orange Mobile Money client
	orangeClient, err := clients.NewMobileMoneyClient(*cfg.GetOrangeConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create Orange client: %w", err)
	}
	manager.SetOrangeClient(orangeClient)

	return manager, nil
}

// setupRouter configures the HTTP router with all routes and middleware
func setupRouter(handler *handlers.ClientHandler) *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Health endpoints
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/ready", handler.ReadinessCheck).Methods("GET")
	router.HandleFunc("/live", handler.LivenessCheck).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Exchange endpoints
	exchange := api.PathPrefix("/exchange").Subrouter()
	
	// Binance endpoints
	binance := exchange.PathPrefix("/binance").Subrouter()
	binance.HandleFunc("/price", handler.GetBinancePrice).Methods("GET")
	binance.HandleFunc("/order", handler.PlaceBinanceOrder).Methods("POST")
	binance.HandleFunc("/order/{orderID}", handler.GetBinanceOrderStatus).Methods("GET")
	binance.HandleFunc("/quotes", handler.GetBinanceQuotes).Methods("GET")
	binance.HandleFunc("/quote/{symbol}", handler.GetBinanceQuote).Methods("GET")

	// Bitget endpoints
	bitget := exchange.PathPrefix("/bitget").Subrouter()
	bitget.HandleFunc("/price", handler.GetBitgetPrice).Methods("GET")
	bitget.HandleFunc("/order", handler.PlaceBitgetOrder).Methods("POST")
	bitget.HandleFunc("/order/{orderID}", handler.GetBitgetOrderStatus).Methods("GET")
	bitget.HandleFunc("/quotes", handler.GetBitgetQuotes).Methods("GET")
	bitget.HandleFunc("/quote/{symbol}", handler.GetBitgetQuote).Methods("GET")

	// Mobile Money endpoints
	mobileMoney := api.PathPrefix("/mobile-money").Subrouter()
	
	// MTN endpoints
	mtn := mobileMoney.PathPrefix("/mtn").Subrouter()
	mtn.HandleFunc("/deposit", handler.ProcessMTNDeposit).Methods("POST")
	mtn.HandleFunc("/status/{transactionID}", handler.GetMTNTransactionStatus).Methods("GET")

	// Orange endpoints
	orange := mobileMoney.PathPrefix("/orange").Subrouter()
	orange.HandleFunc("/deposit", handler.ProcessOrangeDeposit).Methods("POST")
	orange.HandleFunc("/status/{transactionID}", handler.GetOrangeTransactionStatus).Methods("GET")

	// Resilience endpoints
	resilience := api.PathPrefix("/resilience").Subrouter()
	resilience.HandleFunc("/stats", handler.GetResilienceStats).Methods("GET")
	resilience.HandleFunc("/reset", handler.ResetResilienceStats).Methods("POST")

	return router
}
