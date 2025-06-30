package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Mock responses structures
type Quote struct {
	Symbol    string  `json:"symbol"`
	LastPrice string  `json:"lastPrice"`
	AskPrice  string  `json:"askPrice"`
	BidPrice  string  `json:"bidPrice"`
	Volume    string  `json:"volume"`
	Timestamp string  `json:"timestamp"`
}

type QuotesResponse struct {
	Quotes    []Quote `json:"quotes"`
	Timestamp string  `json:"timestamp"`
}

type Order struct {
	OrderID     int64   `json:"orderId"`
	Symbol      string  `json:"symbol"`
	Status      string  `json:"status"`
	Quantity    string  `json:"origQty"`
	Price       string  `json:"price"`
	ExecutedQty string  `json:"executedQty"`
	Timestamp   string  `json:"time"`
}

type OrdersResponse struct {
	Orders    []Order `json:"orders"`
	Timestamp string  `json:"timestamp"`
}

type PriceResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Success   bool    `json:"success"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/api/exchange/binance/v1/quotes", getQuotes)
	e.GET("/api/exchange/binance/v1/orders", getOrders)
	e.POST("/api/workflow/v1/BinancePrice", getBinancePrice)

	log.Println("Starting test server on port :8088")
	e.Logger.Fatal(e.Start(":8088"))
}

func getQuotes(c echo.Context) error {
	quotes := []Quote{
		{
			Symbol:    "BTCUSDT",
			LastPrice: "107838.46",
			AskPrice:  "107840.50",
			BidPrice:  "107835.25",
			Volume:    "1234.56789",
			Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
		{
			Symbol:    "ETHUSDT",
			LastPrice: "3945.67",
			AskPrice:  "3946.12",
			BidPrice:  "3945.23",
			Volume:    "9876.54321",
			Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
		{
			Symbol:    "ADAUSDT",
			LastPrice: "0.8925",
			AskPrice:  "0.8930",
			BidPrice:  "0.8920",
			Volume:    "15678.9876",
			Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
		{
			Symbol:    "DOTUSDT",
			LastPrice: "18.45",
			AskPrice:  "18.47",
			BidPrice:  "18.43",
			Volume:    "5432.1098",
			Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	}

	response := QuotesResponse{
		Quotes:    quotes,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}

func getOrders(c echo.Context) error {
	orders := []Order{
		{
			OrderID:     12345678,
			Symbol:      "BTCUSDT",
			Status:      "FILLED",
			Quantity:    "0.001",
			Price:       "107800.00",
			ExecutedQty: "0.001",
			Timestamp:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
		{
			OrderID:     12345679,
			Symbol:      "ETHUSDT",
			Status:      "PARTIALLY_FILLED",
			Quantity:    "0.5",
			Price:       "3940.00",
			ExecutedQty: "0.25",
			Timestamp:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
		{
			OrderID:     12345680,
			Symbol:      "ADAUSDT",
			Status:      "NEW",
			Quantity:    "100",
			Price:       "0.8900",
			ExecutedQty: "0",
			Timestamp:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	}

	response := OrdersResponse{
		Orders:    orders,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}

func getBinancePrice(c echo.Context) error {
	var symbol string
	if err := json.NewDecoder(c.Request().Body).Decode(&symbol); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid symbol format"})
	}

	// Mock price based on symbol
	var price float64
	switch symbol {
	case "BTCUSDT":
		price = 107838.46
	case "ETHUSDT":
		price = 3945.67
	case "ADAUSDT":
		price = 0.8925
	case "DOTUSDT":
		price = 18.45
	default:
		price = 100.00 // Default price for unknown symbols
	}

	response := PriceResponse{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Format(time.RFC3339),
		Success:   true,
	}

	return c.JSON(http.StatusOK, response)
}
