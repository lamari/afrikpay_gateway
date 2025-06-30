package handler

import (
	"context"

	"github.com/afrikpay/gateway/internal/config"
	"github.com/labstack/echo/v4"
	goTemporalClient "go.temporal.io/sdk/client"
)

func RegisterRoutes(e *echo.Echo, cfg *config.Config) {
	// Endpoint unique pour tous les workflows
	e.POST("/api/workflow/:version/:nameworkflow", WorkflowHandler)
	
	// Binance-specific endpoints
	e.GET("/api/exchange/binance/v1/quotes", BinanceQuotesHandler)
	e.GET("/api/exchange/binance/v1/orders", BinanceOrdersHandler)
	e.POST("/api/exchange/binance/v1/order", BinancePlaceOrderHandler)
	e.GET("/api/exchange/binance/v1/order/:orderId", BinanceGetOrderStatusHandler)
}

type TemporalClientIface interface {
	ExecuteWorkflow(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error)
}

var temporalClient TemporalClientIface

func SetTemporalClient(c TemporalClientIface) {
	temporalClient = c
}
