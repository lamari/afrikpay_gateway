package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BinanceClient structure pour le client API
type BinanceClient struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	HTTPClient *http.Client
}

// AccountInfo structure pour les informations de compte
type AccountInfo struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int64     `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
}

// Balance structure pour les soldes
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// OrderResponse structure pour la réponse d'ordre
type OrderResponse struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Side          string `json:"side"`
}

// APIError structure pour les erreurs API
type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Msg)
}

// NewBinanceClient crée un nouveau client Binance
func NewBinanceClient(apiKey, apiSecret string, testnet bool) *BinanceClient {
	baseURL := "https://api.binance.com"
	if testnet {
		baseURL = "https://testnet.binance.vision"
	}

	return &BinanceClient{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// signRequest signe une requête avec HMAC SHA256
func (c *BinanceClient) signRequest(params url.Values) string {
	queryString := params.Encode()
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// makeRequest effectue une requête HTTP
func (c *BinanceClient) makeRequest(method, endpoint string, params url.Values, signed bool) ([]byte, error) {
	if signed {
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
		params.Set("recvWindow", "5000")
		signature := c.signRequest(params)
		params.Set("signature", signature)
	}

	var req *http.Request
	var err error

	if method == "GET" {
		fullURL := c.BaseURL + endpoint + "?" + params.Encode()
		req, err = http.NewRequest("GET", fullURL, nil)
	} else {
		req, err = http.NewRequest(method, c.BaseURL+endpoint, strings.NewReader(params.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return nil, err
	}

	if signed {
		req.Header.Set("X-MBX-APIKEY", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Vérifier les erreurs API
	if resp.StatusCode != 200 {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetAccountInfo récupère les informations du compte
func (c *BinanceClient) GetAccountInfo() (*AccountInfo, error) {
	params := url.Values{}

	body, err := c.makeRequest("GET", "/api/v3/account", params, true)
	if err != nil {
		return nil, err
	}

	var account AccountInfo
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// PlaceOrder place un ordre
func (c *BinanceClient) PlaceOrder(symbol, side, orderType, quantity, price string) (*OrderResponse, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", orderType)
	params.Set("quantity", quantity)

	if orderType == "LIMIT" && price != "" {
		params.Set("price", price)
		params.Set("timeInForce", "GTC")
	}

	body, err := c.makeRequest("POST", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var order OrderResponse
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderBook récupère le carnet d'ordres
func (c *BinanceClient) GetOrderBook(symbol string, limit int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.makeRequest("GET", "/api/v3/depth", params, false)
	if err != nil {
		return nil, err
	}

	var orderBook map[string]interface{}
	err = json.Unmarshal(body, &orderBook)
	if err != nil {
		return nil, err
	}

	return orderBook, nil
}

// GetPrice récupère le prix actuel d'un symbole
func (c *BinanceClient) GetPrice(symbol string) (string, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	body, err := c.makeRequest("GET", "/api/v3/ticker/price", params, false)
	if err != nil {
		return "", err
	}

	var priceData map[string]interface{}
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		return "", err
	}

	price, ok := priceData["price"].(string)
	if !ok {
		return "", fmt.Errorf("impossible de récupérer le prix")
	}

	return price, nil
}

func main() {
	// Configuration - UTILISEZ VOS VRAIES CLÉS API ICI
	apiKey := "Bjhd5FsunILOwVT1RXJiUxOhZt7MQ2jydKs2jspF11kB6tv1xY3EAmDJvrR4w8la"
	apiSecret := "LIWTukFxAjUH43qx7upKrADY1q8Ogc46OLMbeqgsv98noMvV7mWkGeseKtgUmbtC"

	// Créer le client (testnet = true pour les tests)
	client := NewBinanceClient(apiKey, apiSecret, true)

	// Exemple 1: Récupérer les informations du compte
	fmt.Println("=== Informations du compte ===")
	account, err := client.GetAccountInfo()
	if err != nil {
		fmt.Printf("Erreur lors de la récupération du compte: %v\n", err)
	} else {
		fmt.Printf("Peut trader: %v\n", account.CanTrade)
		fmt.Printf("Type de compte: %s\n", account.AccountType)

		// Afficher quelques soldes non nuls
		fmt.Println("Soldes:")
		for _, balance := range account.Balances {
			if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
				fmt.Printf("  %s: Libre=%s, Bloqué=%s\n",
					balance.Asset, balance.Free, balance.Locked)
			}
		}
	}

	// Exemple 2: Récupérer le prix de BTC/USDT
	fmt.Println("\n=== Prix actuel ===")
	price, err := client.GetPrice("BTCUSDT")
	if err != nil {
		fmt.Printf("Erreur lors de la récupération du prix: %v\n", err)
	} else {
		fmt.Printf("Prix BTC/USDT: %s\n", price)
	}

	// Exemple 3: Récupérer le carnet d'ordres
	fmt.Println("\n=== Carnet d'ordres ===")
	_, err = client.GetOrderBook("BTCUSDT", 5)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération du carnet: %v\n", err)
	} else {
		fmt.Printf("Carnet d'ordres récupéré avec succès\n")
		// Vous pouvez parser et afficher les bids/asks ici
	}

	// Exemple 4: Placer un ordre (ATTENTION: Ceci placera un vrai ordre sur le testnet)
	fmt.Println("\n=== Placement d'ordre (exemple) ===")
	fmt.Println("Décommentez le code ci-dessous pour tester le placement d'ordre")

	order, err := client.PlaceOrder(
		"BTCUSDT", // symbole
		"BUY",     // côté (BUY/SELL)
		"MARKET",  // type d'ordre
		"0.001",   // quantité
		"",        // prix (vide pour market order)
	)
	if err != nil {
		fmt.Printf("Erreur lors du placement de l'ordre: %v\n", err)
	} else {
		fmt.Printf("Ordre placé avec succès!\n")
		fmt.Printf("ID de l'ordre: %d\n", order.OrderID)
		fmt.Printf("Status: %s\n", order.Status)
		fmt.Printf("Quantité exécutée: %s\n", order.ExecutedQty)
	}

}
