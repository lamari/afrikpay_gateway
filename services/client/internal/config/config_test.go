package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_FromFile(t *testing.T) {
	t.Run("Valid YAML config file", func(t *testing.T) {
		// Create temporary config file
		configContent := `
server:
  port: 8004
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s

binance:
  base_url: "https://api.binance.com"
  api_key: "test-binance-key"
  secret_key: "test-binance-secret"
  timeout: 30s
  max_retries: 3

bitget:
  base_url: "https://api.bitget.com"
  api_key: "test-bitget-key"
  secret_key: "test-bitget-secret"
  passphrase: "test-passphrase"
  timeout: 30s
  max_retries: 3

mobile_money:
  mtn:
    base_url: "https://sandbox.momodeveloper.mtn.com"
    api_key: "test-mtn-key"
    api_secret: "test-mtn-secret"
    subscription_key: "test-mtn-subscription"
    timeout: 30s
    max_retries: 3
  orange:
    base_url: "https://api.orange.com"
    api_key: "test-orange-key"
    api_secret: "test-orange-secret"
    subscription_key: "test-orange-subscription"
    timeout: 30s
    max_retries: 3

resilience:
  circuit_breaker:
    failure_threshold: 5
    recovery_timeout: 30s
    timeout: 30s
  retry:
    max_retries: 3
    initial_delay: 100ms
    max_delay: 5s
    jitter: true
  timeout:
    default: 30s
    api: 30s

logging:
  level: "info"
  format: "json"
`

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		require.NoError(t, err)

		config, err := LoadConfig(configPath)

		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify server config
		assert.Equal(t, 8004, config.Server.Port)
		assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)

		// Verify Binance config
		assert.Equal(t, "https://api.binance.com", config.Binance.BaseURL)
		assert.Equal(t, "test-binance-key", config.Binance.APIKey)
		assert.Equal(t, "test-binance-secret", config.Binance.SecretKey)

		// Verify Bitget config
		assert.Equal(t, "https://api.bitget.com", config.Bitget.BaseURL)
		assert.Equal(t, "test-bitget-key", config.Bitget.APIKey)
		assert.Equal(t, "test-bitget-secret", config.Bitget.SecretKey)
		assert.Equal(t, "test-passphrase", config.Bitget.Passphrase)

		// Verify Mobile Money configs
		assert.Equal(t, "https://sandbox.momodeveloper.mtn.com", config.MobileMoney.MTN.BaseURL)
		assert.Equal(t, "test-mtn-key", config.MobileMoney.MTN.APIKey)
		assert.Equal(t, "https://api.orange.com", config.MobileMoney.Orange.BaseURL)
		assert.Equal(t, "test-orange-key", config.MobileMoney.Orange.APIKey)

		// Verify resilience config
		assert.Equal(t, 5, config.Resilience.CircuitBreaker.FailureThreshold)
		assert.Equal(t, 30*time.Second, config.Resilience.CircuitBreaker.RecoveryTimeout)
		assert.Equal(t, 3, config.Resilience.Retry.MaxRetries)
		assert.True(t, config.Resilience.Retry.Jitter)

		// Verify logging config
		assert.Equal(t, "info", config.Logging.Level)
		assert.Equal(t, "json", config.Logging.Format)
	})

	t.Run("Invalid YAML file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
		require.NoError(t, err)

		config, err := LoadConfig(configPath)

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "failed to load config from file")
	})

	t.Run("Non-existent file", func(t *testing.T) {
		config, err := LoadConfig("/non/existent/file.yaml")

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "failed to load config from file")
	})
}

func TestLoadConfig_FromEnvironment(t *testing.T) {
	t.Run("Environment variables override", func(t *testing.T) {
		// Set environment variables
		envVars := map[string]string{
			"CLIENT_PORT":           "9000",
			"BINANCE_BASE_URL":      "https://testnet.binance.vision",
			"BINANCE_API_KEY":       "env-binance-key",
			"BINANCE_SECRET_KEY":    "env-binance-secret",
			"BINANCE_TIMEOUT":       "45s",
			"BINANCE_MAX_RETRIES":   "5",
			"BITGET_BASE_URL":       "https://api.bitget.com",
			"BITGET_API_KEY":        "env-bitget-key",
			"BITGET_SECRET_KEY":     "env-bitget-secret",
			"BITGET_PASSPHRASE":     "env-passphrase",
			"BITGET_TIMEOUT":        "45s",
			"BITGET_MAX_RETRIES":    "5",
			"MTN_BASE_URL":          "https://sandbox.momodeveloper.mtn.com",
			"MTN_API_KEY":           "env-mtn-key",
			"MTN_API_SECRET":        "env-mtn-secret",
			"MTN_SUBSCRIPTION_KEY":  "env-mtn-subscription",
			"MTN_TIMEOUT":           "45s",
			"MTN_MAX_RETRIES":       "5",
			"ORANGE_BASE_URL":       "https://api.orange.com",
			"ORANGE_API_KEY":        "env-orange-key",
			"ORANGE_API_SECRET":     "env-orange-secret",
			"ORANGE_SUBSCRIPTION_KEY": "env-orange-subscription",
			"ORANGE_TIMEOUT":        "45s",
			"ORANGE_MAX_RETRIES":    "5",
			"LOG_LEVEL":             "debug",
			"LOG_FORMAT":            "text",
		}

		// Set environment variables
		for key, value := range envVars {
			t.Setenv(key, value)
		}

		config, err := LoadConfig("")

		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify environment variables were loaded
		assert.Equal(t, 9000, config.Server.Port)
		assert.Equal(t, "https://testnet.binance.vision", config.Binance.BaseURL)
		assert.Equal(t, "env-binance-key", config.Binance.APIKey)
		assert.Equal(t, "env-binance-secret", config.Binance.SecretKey)
		assert.Equal(t, 45*time.Second, config.Binance.Timeout)
		assert.Equal(t, 5, config.Binance.MaxRetries)

		assert.Equal(t, "https://api.bitget.com", config.Bitget.BaseURL)
		assert.Equal(t, "env-bitget-key", config.Bitget.APIKey)
		assert.Equal(t, "env-bitget-secret", config.Bitget.SecretKey)
		assert.Equal(t, "env-passphrase", config.Bitget.Passphrase)
		assert.Equal(t, 45*time.Second, config.Bitget.Timeout)
		assert.Equal(t, 5, config.Bitget.MaxRetries)

		assert.Equal(t, "https://sandbox.momodeveloper.mtn.com", config.MobileMoney.MTN.BaseURL)
		assert.Equal(t, "env-mtn-key", config.MobileMoney.MTN.APIKey)
		assert.Equal(t, "env-mtn-secret", config.MobileMoney.MTN.APISecret)
		assert.Equal(t, "env-mtn-subscription", config.MobileMoney.MTN.SubscriptionKey)
		assert.Equal(t, 45*time.Second, config.MobileMoney.MTN.Timeout)
		assert.Equal(t, 5, config.MobileMoney.MTN.MaxRetries)

		assert.Equal(t, "https://api.orange.com", config.MobileMoney.Orange.BaseURL)
		assert.Equal(t, "env-orange-key", config.MobileMoney.Orange.APIKey)
		assert.Equal(t, "env-orange-secret", config.MobileMoney.Orange.APISecret)
		assert.Equal(t, "env-orange-subscription", config.MobileMoney.Orange.SubscriptionKey)
		assert.Equal(t, 45*time.Second, config.MobileMoney.Orange.Timeout)
		assert.Equal(t, 5, config.MobileMoney.Orange.MaxRetries)

		assert.Equal(t, "debug", config.Logging.Level)
		assert.Equal(t, "text", config.Logging.Format)
	})
}

func TestLoadConfig_Defaults(t *testing.T) {
	t.Run("Default values are set", func(t *testing.T) {
		// Clear any existing environment variables
		envVars := []string{
			"CLIENT_PORT", "BINANCE_BASE_URL", "BINANCE_API_KEY", "BINANCE_SECRET_KEY",
			"BITGET_BASE_URL", "BITGET_API_KEY", "BITGET_SECRET_KEY", "BITGET_PASSPHRASE",
			"MTN_BASE_URL", "MTN_API_KEY", "ORANGE_BASE_URL", "ORANGE_API_KEY",
			"LOG_LEVEL", "LOG_FORMAT",
		}
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		// Set minimal required environment variables to pass validation
		t.Setenv("BINANCE_API_KEY", "test-key")
		t.Setenv("BINANCE_SECRET_KEY", "test-secret")
		t.Setenv("BITGET_API_KEY", "test-key")
		t.Setenv("BITGET_SECRET_KEY", "test-secret")
		t.Setenv("BITGET_PASSPHRASE", "test-passphrase")
		t.Setenv("MTN_API_KEY", "test-key")
		t.Setenv("MTN_API_SECRET", "test-secret")
		t.Setenv("MTN_SUBSCRIPTION_KEY", "test-subscription")
		t.Setenv("ORANGE_API_KEY", "test-key")
		t.Setenv("ORANGE_API_SECRET", "test-secret")
		t.Setenv("ORANGE_SUBSCRIPTION_KEY", "test-subscription")

		config, err := LoadConfig("")

		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Verify default values
		assert.Equal(t, 8004, config.Server.Port)
		assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)
		assert.Equal(t, 30*time.Second, config.Server.WriteTimeout)
		assert.Equal(t, 60*time.Second, config.Server.IdleTimeout)

		assert.Equal(t, "https://api.binance.com", config.Binance.BaseURL)
		assert.Equal(t, 30*time.Second, config.Binance.Timeout)
		assert.Equal(t, 3, config.Binance.MaxRetries)

		assert.Equal(t, "https://api.bitget.com", config.Bitget.BaseURL)
		assert.Equal(t, 30*time.Second, config.Bitget.Timeout)
		assert.Equal(t, 3, config.Bitget.MaxRetries)

		assert.Equal(t, 30*time.Second, config.MobileMoney.MTN.Timeout)
		assert.Equal(t, 3, config.MobileMoney.MTN.MaxRetries)
		assert.Equal(t, 30*time.Second, config.MobileMoney.Orange.Timeout)
		assert.Equal(t, 3, config.MobileMoney.Orange.MaxRetries)

		assert.Equal(t, 5, config.Resilience.CircuitBreaker.FailureThreshold)
		assert.Equal(t, 30*time.Second, config.Resilience.CircuitBreaker.RecoveryTimeout)
		assert.Equal(t, 30*time.Second, config.Resilience.CircuitBreaker.Timeout)
		assert.Equal(t, 3, config.Resilience.Retry.MaxRetries)
		assert.Equal(t, 100*time.Millisecond, config.Resilience.Retry.InitialDelay)
		assert.Equal(t, 5*time.Second, config.Resilience.Retry.MaxDelay)
		assert.Equal(t, 30*time.Second, config.Resilience.Timeout.Default)
		assert.Equal(t, 30*time.Second, config.Resilience.Timeout.API)

		assert.Equal(t, "info", config.Logging.Level)
		assert.Equal(t, "json", config.Logging.Format)
	})
}

func TestConfig_Validate(t *testing.T) {
	t.Run("Valid configuration", func(t *testing.T) {
		config := &Config{
			Server: ServerConfig{
				Port:         8004,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
			},
			Binance: models.BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			Bitget: models.BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			MobileMoney: MobileMoneyConfigs{
				MTN: models.MobileMoneyConfig{
					BaseURL:         "https://sandbox.momodeveloper.mtn.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
				Orange: models.MobileMoneyConfig{
					BaseURL:         "https://api.orange.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
			},
			Resilience: ResilienceConfig{
				CircuitBreaker: CircuitBreakerConfig{
					FailureThreshold: 5,
					RecoveryTimeout:  30 * time.Second,
					Timeout:          30 * time.Second,
				},
				Retry: RetryConfig{
					MaxRetries:   3,
					InitialDelay: 100 * time.Millisecond,
					MaxDelay:     5 * time.Second,
					Jitter:       true,
				},
				Timeout: TimeoutConfig{
					Default: 30 * time.Second,
					API:     30 * time.Second,
				},
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "json",
			},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid server port", func(t *testing.T) {
		config := &Config{
			Server: ServerConfig{Port: 0},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server port")
	})

	t.Run("Invalid log level", func(t *testing.T) {
		config := &Config{
			Server: ServerConfig{Port: 8004},
			Binance: models.BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			Bitget: models.BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			MobileMoney: MobileMoneyConfigs{
				MTN: models.MobileMoneyConfig{
					BaseURL:         "https://sandbox.momodeveloper.mtn.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
				Orange: models.MobileMoneyConfig{
					BaseURL:         "https://api.orange.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
			},
			Resilience: ResilienceConfig{
				CircuitBreaker: CircuitBreakerConfig{
					FailureThreshold: 5,
					RecoveryTimeout:  30 * time.Second,
					Timeout:          30 * time.Second,
				},
				Retry: RetryConfig{
					MaxRetries:   3,
					InitialDelay: 100 * time.Millisecond,
					MaxDelay:     5 * time.Second,
				},
				Timeout: TimeoutConfig{
					Default: 30 * time.Second,
					API:     30 * time.Second,
				},
			},
			Logging: LoggingConfig{
				Level:  "invalid",
				Format: "json",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})

	t.Run("Invalid log format", func(t *testing.T) {
		config := &Config{
			Server: ServerConfig{Port: 8004},
			Binance: models.BinanceConfig{
				BaseURL:    "https://api.binance.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			Bitget: models.BitgetConfig{
				BaseURL:    "https://api.bitget.com",
				APIKey:     "test-key",
				SecretKey:  "test-secret",
				Passphrase: "test-passphrase",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			MobileMoney: MobileMoneyConfigs{
				MTN: models.MobileMoneyConfig{
					BaseURL:         "https://sandbox.momodeveloper.mtn.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
				Orange: models.MobileMoneyConfig{
					BaseURL:         "https://api.orange.com",
					APIKey:          "test-key",
					APISecret:       "test-secret",
					SubscriptionKey: "test-subscription",
					Timeout:         30 * time.Second,
					MaxRetries:      3,
				},
			},
			Resilience: ResilienceConfig{
				CircuitBreaker: CircuitBreakerConfig{
					FailureThreshold: 5,
					RecoveryTimeout:  30 * time.Second,
					Timeout:          30 * time.Second,
				},
				Retry: RetryConfig{
					MaxRetries:   3,
					InitialDelay: 100 * time.Millisecond,
					MaxDelay:     5 * time.Second,
				},
				Timeout: TimeoutConfig{
					Default: 30 * time.Second,
					API:     30 * time.Second,
				},
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "invalid",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log format")
	})
}

func TestConfig_GetterMethods(t *testing.T) {
	config := &Config{
		Binance: models.BinanceConfig{
			BaseURL:   "https://api.binance.com",
			APIKey:    "binance-key",
			SecretKey: "binance-secret",
		},
		Bitget: models.BitgetConfig{
			BaseURL:    "https://api.bitget.com",
			APIKey:     "bitget-key",
			SecretKey:  "bitget-secret",
			Passphrase: "bitget-passphrase",
		},
		MobileMoney: MobileMoneyConfigs{
			MTN: models.MobileMoneyConfig{
				BaseURL:   "https://sandbox.momodeveloper.mtn.com",
				APIKey:    "mtn-key",
				APISecret: "mtn-secret",
			},
			Orange: models.MobileMoneyConfig{
				BaseURL:   "https://api.orange.com",
				APIKey:    "orange-key",
				APISecret: "orange-secret",
			},
		},
	}

	// Test getter methods
	binanceConfig := config.GetBinanceConfig()
	assert.Equal(t, "https://api.binance.com", binanceConfig.BaseURL)
	assert.Equal(t, "binance-key", binanceConfig.APIKey)
	assert.Equal(t, "binance-secret", binanceConfig.SecretKey)

	bitgetConfig := config.GetBitgetConfig()
	assert.Equal(t, "https://api.bitget.com", bitgetConfig.BaseURL)
	assert.Equal(t, "bitget-key", bitgetConfig.APIKey)
	assert.Equal(t, "bitget-secret", bitgetConfig.SecretKey)
	assert.Equal(t, "bitget-passphrase", bitgetConfig.Passphrase)

	mtnConfig := config.GetMTNConfig()
	assert.Equal(t, "https://sandbox.momodeveloper.mtn.com", mtnConfig.BaseURL)
	assert.Equal(t, "mtn-key", mtnConfig.APIKey)
	assert.Equal(t, "mtn-secret", mtnConfig.APISecret)

	orangeConfig := config.GetOrangeConfig()
	assert.Equal(t, "https://api.orange.com", orangeConfig.BaseURL)
	assert.Equal(t, "orange-key", orangeConfig.APIKey)
	assert.Equal(t, "orange-secret", orangeConfig.APISecret)
}
