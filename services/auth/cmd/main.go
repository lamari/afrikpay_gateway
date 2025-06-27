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

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"afrikpay/services/auth/internal/config"
	"afrikpay/services/auth/internal/handlers"
	"afrikpay/services/auth/internal/middleware"
	"afrikpay/services/auth/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT service
	jwtService, err := services.NewJWTService(cfg.JWT.Issuer, "afrikpay-gateway", cfg.JWT.AccessTokenExpiry, cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtService)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	// Setup router
	router := mux.NewRouter()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure appropriately for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Setup middleware
	router.Use(middleware.SimpleLoggingMiddleware)
	router.Use(c.Handler)

	// Health endpoints
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	// Auth endpoints
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")
	authRouter.HandleFunc("/verify", authHandler.Verify).Methods("GET")

	// Protected endpoints (example)
	protectedRouter := router.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware(jwtService))
	protectedRouter.HandleFunc("/profile", authHandler.Profile).Methods("GET")

	// Setup server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Auth service starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
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
