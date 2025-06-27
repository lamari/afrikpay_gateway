package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/clients"
	"github.com/afrikpay/gateway/services/client/internal/models"
	"github.com/gorilla/mux"
)

// ClientHandler handles HTTP requests for the client service
type ClientHandler struct {
	clientManager *clients.ClientManager
}

// NewClientHandler creates a new client handler
func NewClientHandler(clientManager *clients.ClientManager) *ClientHandler {
	return &ClientHandler{
		clientManager: clientManager,
	}
}

// HealthCheck handles health check requests
func (h *ClientHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	health := h.clientManager.HealthCheck(ctx)
	
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"clients":   health,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessCheck handles readiness check requests
func (h *ClientHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	health := h.clientManager.HealthCheck(ctx)
	
	// Check if at least one exchange and one mobile money provider are healthy
	exchangeHealthy := health["binance"] || health["bitget"]
	mobileMoneyHealthy := health["mtn"] || health["orange"]
	
	status := "ready"
	statusCode := http.StatusOK
	
	if !exchangeHealthy || !mobileMoneyHealthy {
		status = "not ready"
		statusCode = http.StatusServiceUnavailable
	}

	response := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"clients":   health,
		"ready":     exchangeHealthy && mobileMoneyHealthy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// LivenessCheck handles liveness check requests
func (h *ClientHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetBinancePrice handles Binance price requests
func (h *ClientHandler) GetBinancePrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBinanceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetPrice(ctx, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PlaceBinanceOrder handles Binance order placement requests
func (h *ClientHandler) PlaceBinanceOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBinanceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.PlaceOrder(ctx, &orderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetBinanceOrderStatus handles Binance order status requests
func (h *ClientHandler) GetBinanceOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	if orderID == "" {
		http.Error(w, "orderID parameter is required", http.StatusBadRequest)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBinanceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetOrderStatus(ctx, symbol, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBinanceQuotes handles Binance quotes requests
func (h *ClientHandler) GetBinanceQuotes(w http.ResponseWriter, r *http.Request) {
	client, err := h.clientManager.GetBinanceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetQuotes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBinanceQuote handles Binance single quote requests
func (h *ClientHandler) GetBinanceQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBinanceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetQuote(ctx, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBitgetPrice handles Bitget price requests
func (h *ClientHandler) GetBitgetPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBitgetClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetPrice(ctx, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PlaceBitgetOrder handles Bitget order placement requests
func (h *ClientHandler) PlaceBitgetOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBitgetClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.PlaceOrder(ctx, &orderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetBitgetOrderStatus handles Bitget order status requests
func (h *ClientHandler) GetBitgetOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	if orderID == "" {
		http.Error(w, "orderID parameter is required", http.StatusBadRequest)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBitgetClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetOrderStatus(ctx, symbol, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBitgetQuotes handles Bitget quotes requests
func (h *ClientHandler) GetBitgetQuotes(w http.ResponseWriter, r *http.Request) {
	client, err := h.clientManager.GetBitgetClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetQuotes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBitgetQuote handles Bitget single quote requests
func (h *ClientHandler) GetBitgetQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetBitgetClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetQuote(ctx, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProcessMTNDeposit handles MTN Mobile Money deposit requests
func (h *ClientHandler) ProcessMTNDeposit(w http.ResponseWriter, r *http.Request) {
	var depositRequest models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DepositRequest to PaymentRequest
	paymentRequest := &models.PaymentRequest{
		Amount:      depositRequest.Amount,
		Currency:    depositRequest.Currency,
		ExternalID:  depositRequest.ExternalID,
		PhoneNumber: depositRequest.PhoneNumber,
		Description: depositRequest.Description,
		CallbackURL: depositRequest.CallbackURL,
		Metadata:    depositRequest.Metadata,
	}

	client, err := h.clientManager.GetMTNClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.InitiatePayment(ctx, paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetMTNTransactionStatus handles MTN transaction status requests
func (h *ClientHandler) GetMTNTransactionStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transactionID"]
	if transactionID == "" {
		http.Error(w, "transactionID parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetMTNClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetPaymentStatus(ctx, transactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProcessOrangeDeposit handles Orange Mobile Money deposit requests
func (h *ClientHandler) ProcessOrangeDeposit(w http.ResponseWriter, r *http.Request) {
	var depositRequest models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DepositRequest to PaymentRequest
	paymentRequest := &models.PaymentRequest{
		Amount:      depositRequest.Amount,
		Currency:    depositRequest.Currency,
		ExternalID:  depositRequest.ExternalID,
		PhoneNumber: depositRequest.PhoneNumber,
		Description: depositRequest.Description,
		CallbackURL: depositRequest.CallbackURL,
		Metadata:    depositRequest.Metadata,
	}

	client, err := h.clientManager.GetOrangeClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.InitiatePayment(ctx, paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetOrangeTransactionStatus handles Orange transaction status requests
func (h *ClientHandler) GetOrangeTransactionStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transactionID"]
	if transactionID == "" {
		http.Error(w, "transactionID parameter is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetOrangeClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	response, err := client.GetPaymentStatus(ctx, transactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetResilienceStats handles resilience statistics requests
func (h *ClientHandler) GetResilienceStats(w http.ResponseWriter, r *http.Request) {
	stats := h.clientManager.GetAllResilienceStats()

	response := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"stats":     stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResetResilienceStats handles resilience statistics reset requests
func (h *ClientHandler) ResetResilienceStats(w http.ResponseWriter, r *http.Request) {
	h.clientManager.ResetAllResilienceStats()

	response := map[string]interface{}{
		"message":   "Resilience statistics reset successfully",
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
