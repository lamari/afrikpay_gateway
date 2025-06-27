package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/clients"
	"github.com/afrikpay/gateway/services/client/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExchangeClient is a mock implementation of ExchangeClient
type MockExchangeClient struct {
	mock.Mock
}

func (m *MockExchangeClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockExchangeClient) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockExchangeClient) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.PriceResponse), args.Error(1)
}

func (m *MockExchangeClient) PlaceOrder(ctx context.Context, order *models.OrderRequest) (*models.OrderResponse, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(*models.OrderResponse), args.Error(1)
}

func (m *MockExchangeClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*models.OrderResponse, error) {
	args := m.Called(ctx, symbol, orderID)
	return args.Get(0).(*models.OrderResponse), args.Error(1)
}

func (m *MockExchangeClient) GetQuotes(ctx context.Context) (*models.QuotesResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.QuotesResponse), args.Error(1)
}

func (m *MockExchangeClient) GetQuote(ctx context.Context, symbol string) (*models.QuoteResponse, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.QuoteResponse), args.Error(1)
}

func (m *MockExchangeClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExchangeClient) GetResilienceStats() *models.ResilienceStats {
	args := m.Called()
	return args.Get(0).(*models.ResilienceStats)
}

func (m *MockExchangeClient) ResetResilienceStats() {
	m.Called()
}

// MockMobileMoneyClient is a mock implementation of MobileMoneyClient
type MockMobileMoneyClient struct {
	mock.Mock
}

func (m *MockMobileMoneyClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMobileMoneyClient) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMobileMoneyClient) InitiatePayment(ctx context.Context, request *models.PaymentRequest) (*models.PaymentResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*models.PaymentResponse), args.Error(1)
}

func (m *MockMobileMoneyClient) GetPaymentStatus(ctx context.Context, referenceID string) (*models.PaymentStatusResponse, error) {
	args := m.Called(ctx, referenceID)
	return args.Get(0).(*models.PaymentStatusResponse), args.Error(1)
}

func (m *MockMobileMoneyClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMobileMoneyClient) GetResilienceStats() *models.ResilienceStats {
	args := m.Called()
	return args.Get(0).(*models.ResilienceStats)
}

func (m *MockMobileMoneyClient) ResetResilienceStats() {
	m.Called()
}

func TestClientHandler_HealthCheck(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// When
	handler.HealthCheck(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "clients")
}

func TestClientHandler_ReadinessCheck_Ready(t *testing.T) {
	// Given
	mockBinance := &MockExchangeClient{}
	mockMTN := &MockMobileMoneyClient{}

	mockBinance.On("HealthCheck", mock.Anything).Return(nil)
	mockMTN.On("HealthCheck", mock.Anything).Return(nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetBinanceClient(mockBinance)
	clientManager.SetMTNClient(mockMTN)

	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	// When
	handler.ReadinessCheck(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ready", response["status"])
	assert.Equal(t, true, response["ready"])
}

func TestClientHandler_ReadinessCheck_NotReady(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	// When
	handler.ReadinessCheck(w, req)

	// Then
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not ready", response["status"])
	assert.Equal(t, false, response["ready"])
}

func TestClientHandler_LivenessCheck(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()

	// When
	handler.LivenessCheck(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alive", response["status"])
	assert.Contains(t, response, "timestamp")
}

func TestClientHandler_GetBinancePrice_Success(t *testing.T) {
	// Given
	mockBinance := &MockExchangeClient{}
	expectedResponse := &models.PriceResponse{
		Symbol:    "BTCUSDT",
		Price:     50000.0,
		Timestamp: time.Now(),
	}

	mockBinance.On("GetPrice", mock.Anything, "BTCUSDT").Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetBinanceClient(mockBinance)

	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/binance/price?symbol=BTCUSDT", nil)
	w := httptest.NewRecorder()

	// When
	handler.GetBinancePrice(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.PriceResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", response.Symbol)
	assert.Equal(t, 50000.0, response.Price)

	mockBinance.AssertExpectations(t)
}

func TestClientHandler_GetBinancePrice_MissingSymbol(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/binance/price", nil)
	w := httptest.NewRecorder()

	// When
	handler.GetBinancePrice(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "symbol parameter is required")
}

func TestClientHandler_PlaceBinanceOrder_Success(t *testing.T) {
	// Given
	mockBinance := &MockExchangeClient{}
	orderRequest := &models.OrderRequest{
		Symbol:   "BTCUSDT",
		Side:     "BUY",
		Type:     "MARKET",
		Quantity: 0.001,
	}

	expectedResponse := &models.OrderResponse{
		OrderID:     "12345",
		Symbol:      "BTCUSDT",
		Status:      "FILLED",
		Side:        "BUY",
		Type:        "MARKET",
		Quantity:    0.001,
		Price:       50000.0,
		ExecutedQty: 0.001,
		Timestamp:   time.Now(),
	}

	mockBinance.On("PlaceOrder", mock.Anything, orderRequest).Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetBinanceClient(mockBinance)

	handler := NewClientHandler(clientManager)

	body, _ := json.Marshal(orderRequest)
	req := httptest.NewRequest("POST", "/binance/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// When
	handler.PlaceBinanceOrder(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "FILLED", response.Status)

	mockBinance.AssertExpectations(t)
}

func TestClientHandler_GetBinanceOrderStatus_Success(t *testing.T) {
	// Given
	mockBinance := &MockExchangeClient{}
	expectedResponse := &models.OrderResponse{
		OrderID:     "12345",
		Symbol:      "BTCUSDT",
		Status:      "FILLED",
		Side:        "BUY",
		Type:        "MARKET",
		Quantity:    0.001,
		Price:       50000.0,
		ExecutedQty: 0.001,
		Timestamp:   time.Now(),
	}

	mockBinance.On("GetOrderStatus", mock.Anything, "BTCUSDT", "12345").Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetBinanceClient(mockBinance)

	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/binance/orders/12345?symbol=BTCUSDT", nil)
	req = mux.SetURLVars(req, map[string]string{"orderID": "12345"})
	w := httptest.NewRecorder()

	// When
	handler.GetBinanceOrderStatus(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "FILLED", response.Status)

	mockBinance.AssertExpectations(t)
}

func TestClientHandler_GetBinanceQuotes_Success(t *testing.T) {
	// Given
	mockBinance := &MockExchangeClient{}
	expectedResponse := &models.QuotesResponse{
		Quotes: []models.QuoteResponse{
			{
				Symbol:    "BTCUSDT",
				BidPrice:  49900.0,
				AskPrice:  50100.0,
				LastPrice: 50000.0,
				Volume:    1000.0,
				Timestamp: time.Now(),
			},
		},
		Timestamp: time.Now(),
	}

	mockBinance.On("GetQuotes", mock.Anything).Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetBinanceClient(mockBinance)

	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/binance/quotes", nil)
	w := httptest.NewRecorder()

	// When
	handler.GetBinanceQuotes(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.QuotesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Quotes, 1)
	assert.Equal(t, "BTCUSDT", response.Quotes[0].Symbol)

	mockBinance.AssertExpectations(t)
}

func TestClientHandler_ProcessMTNDeposit_Success(t *testing.T) {
	// Given
	mockMTN := &MockMobileMoneyClient{}
	depositRequest := &models.DepositRequest{
		Amount:      100.0,
		Currency:    "XAF",
		PhoneNumber: "+237123456789",
		ExternalID:  "ext-123",
		Description: "Test deposit",
	}

	expectedPaymentRequest := &models.PaymentRequest{
		Amount:      100.0,
		Currency:    "XAF",
		PhoneNumber: "+237123456789",
		ExternalID:  "ext-123",
		Description: "Test deposit",
	}

	expectedResponse := &models.PaymentResponse{
		ReferenceID: "ref-123",
		Status:      models.PaymentStatusPending,
		Message:     "Payment initiated",
	}

	mockMTN.On("InitiatePayment", mock.Anything, expectedPaymentRequest).Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetMTNClient(mockMTN)

	handler := NewClientHandler(clientManager)

	body, _ := json.Marshal(depositRequest)
	req := httptest.NewRequest("POST", "/mtn/deposits", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// When
	handler.ProcessMTNDeposit(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ref-123", response.ReferenceID)
	assert.Equal(t, models.PaymentStatusPending, response.Status)

	mockMTN.AssertExpectations(t)
}

func TestClientHandler_GetMTNTransactionStatus_Success(t *testing.T) {
	// Given
	mockMTN := &MockMobileMoneyClient{}
	expectedResponse := &models.PaymentStatusResponse{
		ReferenceID:   "ref-123",
		ExternalID:    "ext-123",
		Status:        models.PaymentStatusSuccess,
		Amount:        100.0,
		Currency:      "XAF",
		PhoneNumber:   "+237123456789",
		TransactionID: "txn-123",
		Message:       "Payment completed",
		Timestamp:     time.Now(),
	}

	mockMTN.On("GetPaymentStatus", mock.Anything, "txn-123").Return(expectedResponse, nil)

	clientManager := &clients.ClientManager{}
	clientManager.SetMTNClient(mockMTN)

	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/mtn/transactions/txn-123", nil)
	req = mux.SetURLVars(req, map[string]string{"transactionID": "txn-123"})
	w := httptest.NewRecorder()

	// When
	handler.GetMTNTransactionStatus(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.PaymentStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ref-123", response.ReferenceID)
	assert.Equal(t, models.PaymentStatusSuccess, response.Status)

	mockMTN.AssertExpectations(t)
}

func TestClientHandler_GetResilienceStats(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	// When
	handler.GetResilienceStats(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "stats")
}

func TestClientHandler_ResetResilienceStats(t *testing.T) {
	// Given
	clientManager := &clients.ClientManager{}
	handler := NewClientHandler(clientManager)

	req := httptest.NewRequest("POST", "/stats/reset", nil)
	w := httptest.NewRecorder()

	// When
	handler.ResetResilienceStats(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "timestamp")
	assert.Equal(t, "Resilience statistics reset successfully", response["message"])
}
