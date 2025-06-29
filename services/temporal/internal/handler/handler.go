package handler

import (
	"net/http"

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
	case "BinanceExchangeRate_v1":
		// TODO
	case "BinanceDeposit_v1":
		// TODO
	case "BinanceBuyCrypto_v1":
		// TODO
	case "BinanceErrorRecovery_v1":
		// TODO
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
