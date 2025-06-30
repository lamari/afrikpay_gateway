package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/afrikpay/gateway/internal/models"
	"github.com/afrikpay/gateway/internal/workflows"
	"github.com/labstack/echo/v4"
	goTemporalClient "go.temporal.io/sdk/client"
)

// Handler unique pour tous les workflows
func WorkflowHandler(c echo.Context) error {
	workflowName := c.Param("nameworkflow")
	version := c.Param("version")
	workflowKey := workflowName + "_" + version
	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}
	var result interface{}

	switch workflowKey {
	// Nouveaux workflows Binance
	case "BinancePrice_v1":
		// Parse input: symbol from request body
		var symbol string
		if err := c.Bind(&symbol); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
		}

		// Execute workflow with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		wf, err := temporalClient.ExecuteWorkflow(
			ctx,
			workflowOptions,
			workflows.BinancePriceWorkflow,
			symbol,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
		}

		// Wait for result with timeout
		var priceResponse interface{}
		err = wf.Get(ctx, &priceResponse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
		}

		result = priceResponse

	case "BinanceQuotes_v1":
		// Execute BinanceQuotes workflow (no input required)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		wf, err := temporalClient.ExecuteWorkflow(
			ctx,
			workflowOptions,
			workflows.BinanceQuotesWorkflow,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
		}

		// Wait for result with timeout
		var quotesResponse interface{}
		err = wf.Get(ctx, &quotesResponse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
		}

		result = quotesResponse

	case "BinanceOrders_v1":
		// Execute BinanceOrders workflow (no input required)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		wf, err := temporalClient.ExecuteWorkflow(
			ctx,
			workflowOptions,
			workflows.BinanceOrdersWorkflow,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
		}

		// Wait for result with timeout
		var ordersResponse interface{}
		err = wf.Get(ctx, &ordersResponse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
		}

		result = ordersResponse

	case "MTNPayment_v1":
		// Parse input: paymentRequest from request body
		var paymentRequest models.PaymentRequest
		if err := c.Bind(&paymentRequest); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
		}

		// Execute MTNPayment workflow
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		wf, err := temporalClient.ExecuteWorkflow(
			ctx,
			workflowOptions,
			workflows.MTNPaymentWorkflow,
			&paymentRequest,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
		}

		// Wait for result
		var paymentResponse interface{}
		err = wf.Get(ctx, &paymentResponse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
		}

		result = paymentResponse

	default:
		return c.JSON(http.StatusNotFound, map[string]string{"error": "workflow not found"})
	}

	return c.JSON(http.StatusOK, result)
}

// BinanceQuotesHandler handles GET requests for Binance quotes
func BinanceQuotesHandler(c echo.Context) error {
	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}

	// Execute BinanceQuotes workflow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wf, err := temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.BinanceQuotesWorkflow,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
	}

	// Wait for result with timeout
	var quotesResponse interface{}
	err = wf.Get(ctx, &quotesResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, quotesResponse)
}

// BinanceOrdersHandler handles GET requests for Binance orders
func BinanceOrdersHandler(c echo.Context) error {
	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}

	// Execute BinanceOrders workflow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wf, err := temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.BinanceOrdersWorkflow,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
	}

	// Wait for result with timeout
	var ordersResponse interface{}
	err = wf.Get(ctx, &ordersResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, ordersResponse)
}

// BinancePlaceOrderHandler handles POST requests for placing orders
func BinancePlaceOrderHandler(c echo.Context) error {
	// Parse request body
	var orderRequest models.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order request: "+err.Error())
	}

	// Validate the order request
	if err := orderRequest.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order: "+err.Error())
	}

	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}

	// Execute BinancePlaceOrder workflow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wf, err := temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.BinancePlaceOrderWorkflow,
		&orderRequest,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
	}

	// Wait for result with timeout
	var orderResponse models.OrderResponse
	err = wf.Get(ctx, &orderResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, orderResponse)
}

// BinanceGetOrderStatusHandler handles GET requests for order status
func BinanceGetOrderStatusHandler(c echo.Context) error {
	// Get path parameters
	orderID := c.Param("orderId")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Order ID is required")
	}

	// Get symbol from query parameter
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Symbol is required as query parameter")
	}

	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}

	// Execute BinanceGetOrderStatus workflow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wf, err := temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.BinanceGetOrderStatusWorkflow,
		symbol,
		orderID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
	}

	// Wait for result with timeout
	var orderResponse models.OrderResponse
	err = wf.Get(ctx, &orderResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, orderResponse)
}

// MTNPaymentHandler handles POST requests for initiating MTN payments
func MTNPaymentHandler(c echo.Context) error {
	// Parse request body
	var paymentRequest models.PaymentRequest
	if err := c.Bind(&paymentRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment request: "+err.Error())
	}

	workflowOptions := goTemporalClient.StartWorkflowOptions{
		TaskQueue: "afrikpay",
	}

	// Execute MTNPaymentWorkflow with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wf, err := temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.MTNPaymentWorkflow,
		&paymentRequest,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start workflow: "+err.Error())
	}

	// Wait for workflow result
	var paymentResponse models.PaymentResponse
	if err := wf.Get(ctx, &paymentResponse); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Workflow failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, paymentResponse)
}
