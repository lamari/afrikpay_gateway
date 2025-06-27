package clients

import (
	"context"
	"errors"
	"testing"

	"github.com/afrikpay/gateway/services/client/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExchangeClient is a mock implementation of ExchangeClient
type MockExchangeClient struct {
	mock.Mock
}

func (m *MockExchangeClient) GetPrice(ctx context.Context, symbol string) (*models.PriceResponse, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.PriceResponse), args.Error(1)
}

func (m *MockExchangeClient) PlaceOrder(ctx context.Context, request *models.OrderRequest) (*models.OrderResponse, error) {
	args := m.Called(ctx, request)
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

func (m *MockExchangeClient) GetResilienceStats() *models.ResilienceStats {
	args := m.Called()
	return args.Get(0).(*models.ResilienceStats)
}

func (m *MockExchangeClient) ResetResilienceStats() {
	m.Called()
}

func (m *MockExchangeClient) Close() error {
	args := m.Called()
	return args.Error(0)
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

func (m *MockMobileMoneyClient) GetResilienceStats() *models.ResilienceStats {
	args := m.Called()
	return args.Get(0).(*models.ResilienceStats)
}

func (m *MockMobileMoneyClient) ResetResilienceStats() {
	m.Called()
}

func (m *MockMobileMoneyClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewClientManager(t *testing.T) {
	manager := NewClientManager()
	assert.NotNil(t, manager)
}

func TestClientManager_SetAndGetClients(t *testing.T) {
	manager := NewClientManager()

	// Test setting and getting Binance client
	mockBinance := &MockExchangeClient{}
	manager.SetBinanceClient(mockBinance)

	binanceClient, err := manager.GetBinanceClient()
	assert.NoError(t, err)
	assert.Equal(t, mockBinance, binanceClient)

	// Test setting and getting Bitget client
	mockBitget := &MockExchangeClient{}
	manager.SetBitgetClient(mockBitget)

	bitgetClient, err := manager.GetBitgetClient()
	assert.NoError(t, err)
	assert.Equal(t, mockBitget, bitgetClient)

	// Test setting and getting MTN client
	mockMTN := &MockMobileMoneyClient{}
	manager.SetMTNClient(mockMTN)

	mtnClient, err := manager.GetMTNClient()
	assert.NoError(t, err)
	assert.Equal(t, mockMTN, mtnClient)

	// Test setting and getting Orange client
	mockOrange := &MockMobileMoneyClient{}
	manager.SetOrangeClient(mockOrange)

	orangeClient, err := manager.GetOrangeClient()
	assert.NoError(t, err)
	assert.Equal(t, mockOrange, orangeClient)
}

func TestClientManager_GetClientErrors(t *testing.T) {
	manager := NewClientManager()

	// Test getting uninitialized clients
	_, err := manager.GetBinanceClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Binance client not initialized")

	_, err = manager.GetBitgetClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Bitget client not initialized")

	_, err = manager.GetMTNClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MTN client not initialized")

	_, err = manager.GetOrangeClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Orange client not initialized")
}

func TestClientManager_GetExchangeClient(t *testing.T) {
	manager := NewClientManager()

	// Test with uninitialized clients
	_, err := manager.GetExchangeClient("binance")
	assert.Error(t, err)

	_, err = manager.GetExchangeClient("bitget")
	assert.Error(t, err)

	// Test with unsupported exchange
	_, err = manager.GetExchangeClient("unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported exchange")

	// Test with initialized clients
	mockBinance := &MockExchangeClient{}
	manager.SetBinanceClient(mockBinance)

	client, err := manager.GetExchangeClient("binance")
	assert.NoError(t, err)
	assert.Equal(t, mockBinance, client)

	mockBitget := &MockExchangeClient{}
	manager.SetBitgetClient(mockBitget)

	client, err = manager.GetExchangeClient("bitget")
	assert.NoError(t, err)
	assert.Equal(t, mockBitget, client)
}

func TestClientManager_GetMobileMoneyClient(t *testing.T) {
	manager := NewClientManager()

	// Test with uninitialized clients
	_, err := manager.GetMobileMoneyClient("mtn")
	assert.Error(t, err)

	_, err = manager.GetMobileMoneyClient("orange")
	assert.Error(t, err)

	// Test with unsupported provider
	_, err = manager.GetMobileMoneyClient("unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported mobile money provider")

	// Test with initialized clients
	mockMTN := &MockMobileMoneyClient{}
	manager.SetMTNClient(mockMTN)

	client, err := manager.GetMobileMoneyClient("mtn")
	assert.NoError(t, err)
	assert.Equal(t, mockMTN, client)

	mockOrange := &MockMobileMoneyClient{}
	manager.SetOrangeClient(mockOrange)

	client, err = manager.GetMobileMoneyClient("orange")
	assert.NoError(t, err)
	assert.Equal(t, mockOrange, client)
}

func TestClientManager_GetAllResilienceStats(t *testing.T) {
	manager := NewClientManager()

	// Test with no clients
	stats := manager.GetAllResilienceStats()
	assert.Empty(t, stats)

	// Test with clients
	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}
	mockMTN := &MockMobileMoneyClient{}
	mockOrange := &MockMobileMoneyClient{}

	binanceStats := &models.ResilienceStats{
		CircuitBreakerState: "closed",
		TotalRequests:       100,
		SuccessfulRequests:  95,
		FailedRequests:      5,
	}
	bitgetStats := &models.ResilienceStats{
		CircuitBreakerState: "closed",
		TotalRequests:       50,
		SuccessfulRequests:  48,
		FailedRequests:      2,
	}
	mtnStats := &models.ResilienceStats{
		CircuitBreakerState: "closed",
		TotalRequests:       20,
		SuccessfulRequests:  20,
		FailedRequests:      0,
	}
	orangeStats := &models.ResilienceStats{
		CircuitBreakerState: "closed",
		TotalRequests:       15,
		SuccessfulRequests:  14,
		FailedRequests:      1,
	}

	mockBinance.On("GetResilienceStats").Return(binanceStats)
	mockBitget.On("GetResilienceStats").Return(bitgetStats)
	mockMTN.On("GetResilienceStats").Return(mtnStats)
	mockOrange.On("GetResilienceStats").Return(orangeStats)

	manager.SetBinanceClient(mockBinance)
	manager.SetBitgetClient(mockBitget)
	manager.SetMTNClient(mockMTN)
	manager.SetOrangeClient(mockOrange)

	stats = manager.GetAllResilienceStats()

	assert.Len(t, stats, 4)
	assert.Equal(t, binanceStats, stats["binance"])
	assert.Equal(t, bitgetStats, stats["bitget"])
	assert.Equal(t, mtnStats, stats["mtn"])
	assert.Equal(t, orangeStats, stats["orange"])

	mockBinance.AssertExpectations(t)
	mockBitget.AssertExpectations(t)
	mockMTN.AssertExpectations(t)
	mockOrange.AssertExpectations(t)
}

func TestClientManager_ResetAllResilienceStats(t *testing.T) {
	manager := NewClientManager()

	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}
	mockMTN := &MockMobileMoneyClient{}
	mockOrange := &MockMobileMoneyClient{}

	mockBinance.On("ResetResilienceStats").Return()
	mockBitget.On("ResetResilienceStats").Return()
	mockMTN.On("ResetResilienceStats").Return()
	mockOrange.On("ResetResilienceStats").Return()

	manager.SetBinanceClient(mockBinance)
	manager.SetBitgetClient(mockBitget)
	manager.SetMTNClient(mockMTN)
	manager.SetOrangeClient(mockOrange)

	manager.ResetAllResilienceStats()

	mockBinance.AssertExpectations(t)
	mockBitget.AssertExpectations(t)
	mockMTN.AssertExpectations(t)
	mockOrange.AssertExpectations(t)
}

func TestClientManager_HealthCheck(t *testing.T) {
	manager := NewClientManager()
	ctx := context.Background()

	// Test with no clients
	health := manager.HealthCheck(ctx)
	assert.Equal(t, false, health["binance"])
	assert.Equal(t, false, health["bitget"])
	assert.Equal(t, false, health["mtn"])
	assert.Equal(t, false, health["orange"])

	// Test with clients
	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}
	mockMTN := &MockMobileMoneyClient{}
	mockOrange := &MockMobileMoneyClient{}

	// Mock successful health checks for exchanges
	mockBinance.On("GetQuotes", ctx).Return(&models.QuotesResponse{}, nil)
	mockBitget.On("GetQuotes", ctx).Return((*models.QuotesResponse)(nil), errors.New("connection error"))

	manager.SetBinanceClient(mockBinance)
	manager.SetBitgetClient(mockBitget)
	manager.SetMTNClient(mockMTN)
	manager.SetOrangeClient(mockOrange)

	health = manager.HealthCheck(ctx)

	assert.Equal(t, true, health["binance"])   // Successful GetQuotes call
	assert.Equal(t, false, health["bitget"])   // Failed GetQuotes call
	assert.Equal(t, true, health["mtn"])       // Client is initialized
	assert.Equal(t, true, health["orange"])    // Client is initialized

	mockBinance.AssertExpectations(t)
	mockBitget.AssertExpectations(t)
}

func TestClientManager_Close(t *testing.T) {
	manager := NewClientManager()

	// Test with no clients
	err := manager.Close()
	assert.NoError(t, err)

	// Test with clients that close successfully
	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}
	mockMTN := &MockMobileMoneyClient{}
	mockOrange := &MockMobileMoneyClient{}

	mockBinance.On("Close").Return(nil)
	mockBitget.On("Close").Return(nil)
	mockMTN.On("Close").Return(nil)
	mockOrange.On("Close").Return(nil)

	manager.SetBinanceClient(mockBinance)
	manager.SetBitgetClient(mockBitget)
	manager.SetMTNClient(mockMTN)
	manager.SetOrangeClient(mockOrange)

	err = manager.Close()
	assert.NoError(t, err)

	mockBinance.AssertExpectations(t)
	mockBitget.AssertExpectations(t)
	mockMTN.AssertExpectations(t)
	mockOrange.AssertExpectations(t)
}

func TestClientManager_CloseWithErrors(t *testing.T) {
	manager := NewClientManager()

	// Test with clients that fail to close
	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}

	mockBinance.On("Close").Return(errors.New("binance close error"))
	mockBitget.On("Close").Return(errors.New("bitget close error"))

	manager.SetBinanceClient(mockBinance)
	manager.SetBitgetClient(mockBitget)

	err := manager.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "errors closing clients")
	assert.Contains(t, err.Error(), "binance close error")
	assert.Contains(t, err.Error(), "bitget close error")

	mockBinance.AssertExpectations(t)
	mockBitget.AssertExpectations(t)
}

func TestClientManager_GetSupportedExchanges(t *testing.T) {
	manager := NewClientManager()

	// Test with no clients
	exchanges := manager.GetSupportedExchanges()
	assert.Empty(t, exchanges)

	// Test with clients
	mockBinance := &MockExchangeClient{}
	mockBitget := &MockExchangeClient{}

	manager.SetBinanceClient(mockBinance)
	exchanges = manager.GetSupportedExchanges()
	assert.Contains(t, exchanges, "binance")
	assert.Len(t, exchanges, 1)

	manager.SetBitgetClient(mockBitget)
	exchanges = manager.GetSupportedExchanges()
	assert.Contains(t, exchanges, "binance")
	assert.Contains(t, exchanges, "bitget")
	assert.Len(t, exchanges, 2)
}

func TestClientManager_GetSupportedMobileMoneyProviders(t *testing.T) {
	manager := NewClientManager()

	// Test with no clients
	providers := manager.GetSupportedMobileMoneyProviders()
	assert.Empty(t, providers)

	// Test with clients
	mockMTN := &MockMobileMoneyClient{}
	mockOrange := &MockMobileMoneyClient{}

	manager.SetMTNClient(mockMTN)
	providers = manager.GetSupportedMobileMoneyProviders()
	assert.Contains(t, providers, "mtn")
	assert.Len(t, providers, 1)

	manager.SetOrangeClient(mockOrange)
	providers = manager.GetSupportedMobileMoneyProviders()
	assert.Contains(t, providers, "mtn")
	assert.Contains(t, providers, "orange")
	assert.Len(t, providers, 2)
}
