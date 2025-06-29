package handler

import (
	"context"
	"net/http"
	"time"

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
	var (
		result interface{}
		err    error
	)

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

	default:
		return c.JSON(http.StatusNotFound, map[string]string{"error": "workflow not found"})
	}
	if err != nil {
		// If the error is due to input binding (invalid JSON), return HTTP 400
		if he, ok := err.(*echo.HTTPError); ok {
			return he
		}
		if err == echo.ErrUnsupportedMediaType || err == echo.ErrBadRequest || err.Error() == "code=400, message=Unmarshal type error" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
