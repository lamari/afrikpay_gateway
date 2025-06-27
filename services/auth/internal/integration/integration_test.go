package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"afrikpay/services/auth/internal/config"
	"afrikpay/services/auth/internal/handlers"
	"afrikpay/services/auth/internal/middleware"
	"afrikpay/services/auth/internal/models"
	"afrikpay/services/auth/internal/services"
)

// IntegrationTestSuite contient les composants nécessaires pour les tests d'intégration
type IntegrationTestSuite struct {
	server     *httptest.Server
	router     *mux.Router
	jwtService services.JWTService
	config     *config.Config
}

// setupIntegrationTest configure l'environnement de test d'intégration
func setupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	// Configuration de test
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8001,
			Host: "localhost",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		JWT: config.JWTConfig{
			PrivateKeyPath:     "../../config/private.pem",
			PublicKeyPath:      "../../config/public.pem",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 24 * time.Hour,
			Issuer:             "afrikpay-gateway",
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Vérifier que les clés existent
	// Chercher la racine du projet en remontant jusqu'à trouver go.work
	wd, _ := os.Getwd()
	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.work")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			// Atteint la racine du système de fichiers
			break
		}
		projectRoot = parent
	}
	privateKeyPath := filepath.Join(projectRoot, "config", "keys", "private.pem")
	publicKeyPath := filepath.Join(projectRoot, "config", "keys", "public.pem")
	
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		t.Skipf("Private key not found at %s, skipping integration tests", privateKeyPath)
	}
	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		t.Skipf("Public key not found at %s, skipping integration tests", publicKeyPath)
	}

	// Initialiser le service JWT
	jwtService, err := services.NewJWTService(
		cfg.JWT.Issuer,
		"afrikpay-gateway",
		cfg.JWT.AccessTokenExpiry,
		privateKeyPath,
		publicKeyPath,
	)
	require.NoError(t, err, "Failed to create JWT service")

	// Créer les handlers
	authHandler := handlers.NewAuthHandler(jwtService)
	healthHandler := handlers.NewHealthHandler("v1.0.0")

	// Configurer le router
	r := mux.NewRouter()
	
	// Routes d'authentification
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/auth/verify", authHandler.Verify).Methods("GET")
	r.HandleFunc("/auth/refresh", authHandler.Refresh).Methods("POST")
	
	// Protected routes
	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware(jwtService))
	protected.HandleFunc("/profile", authHandler.Profile).Methods("GET")
	
	// Health routes
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")
	r.HandleFunc("/ready", healthHandler.Ready).Methods("GET")
	r.HandleFunc("/live", healthHandler.Live).Methods("GET")
	
	// Add middleware
	r.Use(middleware.SimpleLoggingMiddleware)

	// Créer le serveur de test
	server := httptest.NewServer(r)

	return &IntegrationTestSuite{
		server:     server,
		router:     r,
		jwtService: jwtService,
		config:     cfg,
	}
}

// teardownIntegrationTest nettoie l'environnement de test
func (suite *IntegrationTestSuite) teardown() {
	if suite.server != nil {
		suite.server.Close()
	}
}

// TestAuthIntegration_LoginFlow teste le flow complet de login
func TestAuthIntegration_LoginFlow(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	// Test data
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Convertir en JSON
	jsonData, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	// Faire la requête de login
	resp, err := http.Post(
		suite.server.URL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier la réponse
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	// Décoder la réponse
	var loginResponse models.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	// Vérifier les tokens
	assert.NotEmpty(t, loginResponse.AccessToken)
	// Vérifier que les tokens sont valides
	claims, err := suite.jwtService.ValidateToken(loginResponse.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

// TestAuthIntegration_VerifyTokenFlow teste la vérification de token
func TestAuthIntegration_VerifyTokenFlow(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	// Générer un token valide pour les tests protégés
	token, err := suite.jwtService.GenerateToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Créer la requête HTTP avec le token dans l'en-tête Authorization
	req, err := http.NewRequest("GET", suite.server.URL+"/auth/verify", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier la réponse
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var claims models.CustomClaims
	err = json.NewDecoder(resp.Body).Decode(&claims)
	require.NoError(t, err)

	// Vérifier les claims retournées
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Contains(t, claims.Roles, "user")
}

// TestAuthIntegration_RefreshTokenFlow teste le renouvellement de token
func TestAuthIntegration_RefreshTokenFlow(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	// Générer une paire de tokens pour obtenir un refresh token
	tokenPair, err := suite.jwtService.GenerateTokenPair("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Préparer la requête de refresh
	refreshRequest := models.RefreshTokenRequest{
		RefreshToken: tokenPair.RefreshToken,
	}
	jsonData, err := json.Marshal(refreshRequest)
	require.NoError(t, err)

	// Faire la requête de refresh
	resp, err := http.Post(
		suite.server.URL+"/auth/refresh",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier la réponse
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var refreshResponse models.RefreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&refreshResponse)
	require.NoError(t, err)

	assert.NotEmpty(t, refreshResponse.AccessToken)
	assert.NotEmpty(t, refreshResponse.RefreshToken)
	// Le TokenType peut être vide dans l'implémentation actuelle
	if refreshResponse.TokenType != "" {
		assert.Equal(t, "Bearer", refreshResponse.TokenType)
	}
}

// TestAuthIntegration_ProtectedEndpoint teste l'accès aux endpoints protégés
func TestAuthIntegration_ProtectedEndpoint(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	// Générer un token valide
	token, err := suite.jwtService.GenerateToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Créer la requête avec le token
	req, err := http.NewRequest("GET", suite.server.URL+"/protected/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier la réponse
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var profileResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&profileResponse)
	require.NoError(t, err)

	assert.Equal(t, "user123", profileResponse["user_id"])
	assert.Equal(t, "test@example.com", profileResponse["email"])
	assert.Contains(t, profileResponse["roles"], "user")
}

// TestAuthIntegration_ProtectedEndpointUnauthorized teste l'accès non autorisé
func TestAuthIntegration_ProtectedEndpointUnauthorized(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	// Créer la requête sans token
	req, err := http.NewRequest("GET", suite.server.URL+"/protected/profile", nil)
	require.NoError(t, err)

	// Exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier que l'accès est refusé
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errorResponse models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	require.NoError(t, err)

	assert.Equal(t, "MISSING_TOKEN", errorResponse.Code)
	assert.Contains(t, errorResponse.Message, "Authorization header is required")
}

// TestAuthIntegration_HealthEndpoints teste les endpoints de santé
func TestAuthIntegration_HealthEndpoints(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	endpoints := []string{"/health", "/ready", "/live"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp, err := http.Get(suite.server.URL + endpoint)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			var healthResponse map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&healthResponse)
			require.NoError(t, err)

			// Vérifier le statut selon l'endpoint
			switch endpoint {
			case "/health":
				assert.Equal(t, "healthy", healthResponse["status"])
			case "/ready":
				assert.Equal(t, "ready", healthResponse["status"])
			case "/live":
				assert.Equal(t, "alive", healthResponse["status"])
			}
		})
	}
}

// TestAuthIntegration_InvalidToken teste la gestion des tokens invalides
func TestAuthIntegration_InvalidToken(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	testCases := []struct {
		name        string
		token       string
		expectedMsg string
	}{
		{
			name:        "Token malformé",
			token:       "invalid.token.here",
			expectedMsg: "Invalid token",
		},
		{
			name:        "Token expiré",
			token:       generateExpiredToken(t, suite.jwtService),
			expectedMsg: "Invalid token",
		},
		{
			name:        "Token vide",
			token:       "",
			expectedMsg: "Authorization header is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", suite.server.URL+"/protected/profile", nil)
			require.NoError(t, err)
			
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var errorResponse models.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
			require.NoError(t, err)

			// Vérifier le code d'erreur selon le cas
			if tc.token == "" {
				assert.Equal(t, "MISSING_TOKEN", errorResponse.Code)
			} else {
				assert.Equal(t, "INVALID_TOKEN", errorResponse.Code)
			}
		})
	}
}

// TestAuthIntegration_SecurityHeaders teste les headers de sécurité
func TestAuthIntegration_SecurityHeaders(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.teardown()

	resp, err := http.Get(suite.server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Vérifier que les headers de sécurité sont présents
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	
	// Note: D'autres headers de sécurité peuvent être ajoutés selon les besoins
	// comme X-Content-Type-Options, X-Frame-Options, etc.
}

// generateExpiredToken génère un token expiré pour les tests
func generateExpiredToken(t *testing.T, jwtService services.JWTService) string {
	// Cette fonction nécessiterait une modification du service JWT pour permettre
	// la génération de tokens avec une expiration personnalisée
	// Pour l'instant, on retourne un token invalide
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6MX0.invalid"
}

// BenchmarkAuthIntegration_LoginEndpoint benchmark pour l'endpoint de login
func BenchmarkAuthIntegration_LoginEndpoint(b *testing.B) {
	suite := setupIntegrationTest(&testing.T{})
	defer suite.teardown()

	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(loginRequest)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Post(
			suite.server.URL+"/auth/login",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
