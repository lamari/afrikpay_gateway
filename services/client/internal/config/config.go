package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/afrikpay/gateway/services/client/internal/models"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the client service
type Config struct {
	Server       ServerConfig                    `yaml:"server"`
	Binance      models.BinanceConfig           `yaml:"binance"`
	Bitget       models.BitgetConfig            `yaml:"bitget"`
	MobileMoney  MobileMoneyConfigs             `yaml:"mobile_money"`
	Resilience   ResilienceConfig               `yaml:"resilience"`
	Logging      LoggingConfig                  `yaml:"logging"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// MobileMoneyConfigs holds configurations for mobile money providers
type MobileMoneyConfigs struct {
	MTN    models.MobileMoneyConfig `yaml:"mtn"`
	Orange models.MobileMoneyConfig `yaml:"orange"`
}

// ResilienceConfig holds resilience pattern configurations
type ResilienceConfig struct {
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
	Retry          RetryConfig          `yaml:"retry"`
	Timeout        TimeoutConfig        `yaml:"timeout"`
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int           `yaml:"failure_threshold"`
	RecoveryTimeout  time.Duration `yaml:"recovery_timeout"`
	Timeout          time.Duration `yaml:"timeout"`
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries   int           `yaml:"max_retries"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
	Jitter       bool          `yaml:"jitter"`
}

// TimeoutConfig holds timeout configuration
type TimeoutConfig struct {
	Default time.Duration `yaml:"default"`
	API     time.Duration `yaml:"api"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Load from YAML file if provided
	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	if err := loadFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// Set defaults
	setDefaults(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from YAML file
func loadFromFile(config *Config, configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) error {
	// Server configuration
	if port := os.Getenv("CLIENT_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// Binance configuration
	if baseURL := os.Getenv("BINANCE_BASE_URL"); baseURL != "" {
		config.Binance.BaseURL = baseURL
	}
	if apiKey := os.Getenv("BINANCE_API_KEY"); apiKey != "" {
		config.Binance.APIKey = apiKey
	}
	if secretKey := os.Getenv("BINANCE_SECRET_KEY"); secretKey != "" {
		config.Binance.SecretKey = secretKey
	}
	if timeout := os.Getenv("BINANCE_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.Binance.Timeout = t
		}
	}
	if maxRetries := os.Getenv("BINANCE_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil {
			config.Binance.MaxRetries = r
		}
	}

	// Bitget configuration
	if baseURL := os.Getenv("BITGET_BASE_URL"); baseURL != "" {
		config.Bitget.BaseURL = baseURL
	}
	if apiKey := os.Getenv("BITGET_API_KEY"); apiKey != "" {
		config.Bitget.APIKey = apiKey
	}
	if secretKey := os.Getenv("BITGET_SECRET_KEY"); secretKey != "" {
		config.Bitget.SecretKey = secretKey
	}
	if passphrase := os.Getenv("BITGET_PASSPHRASE"); passphrase != "" {
		config.Bitget.Passphrase = passphrase
	}
	if timeout := os.Getenv("BITGET_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.Bitget.Timeout = t
		}
	}
	if maxRetries := os.Getenv("BITGET_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil {
			config.Bitget.MaxRetries = r
		}
	}

	// MTN Mobile Money configuration
	if baseURL := os.Getenv("MTN_BASE_URL"); baseURL != "" {
		config.MobileMoney.MTN.BaseURL = baseURL
	}
	if apiKey := os.Getenv("MTN_API_KEY"); apiKey != "" {
		config.MobileMoney.MTN.APIKey = apiKey
	}
	if apiSecret := os.Getenv("MTN_API_SECRET"); apiSecret != "" {
		config.MobileMoney.MTN.APISecret = apiSecret
	}
	if subscriptionKey := os.Getenv("MTN_SUBSCRIPTION_KEY"); subscriptionKey != "" {
		config.MobileMoney.MTN.SubscriptionKey = subscriptionKey
	}
	if timeout := os.Getenv("MTN_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.MobileMoney.MTN.Timeout = t
		}
	}
	if maxRetries := os.Getenv("MTN_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil {
			config.MobileMoney.MTN.MaxRetries = r
		}
	}

	// Orange Mobile Money configuration
	if baseURL := os.Getenv("ORANGE_BASE_URL"); baseURL != "" {
		config.MobileMoney.Orange.BaseURL = baseURL
	}
	if apiKey := os.Getenv("ORANGE_API_KEY"); apiKey != "" {
		config.MobileMoney.Orange.APIKey = apiKey
	}
	if apiSecret := os.Getenv("ORANGE_API_SECRET"); apiSecret != "" {
		config.MobileMoney.Orange.APISecret = apiSecret
	}
	if subscriptionKey := os.Getenv("ORANGE_SUBSCRIPTION_KEY"); subscriptionKey != "" {
		config.MobileMoney.Orange.SubscriptionKey = subscriptionKey
	}
	if timeout := os.Getenv("ORANGE_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.MobileMoney.Orange.Timeout = t
		}
	}
	if maxRetries := os.Getenv("ORANGE_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil {
			config.MobileMoney.Orange.MaxRetries = r
		}
	}

	// Logging configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	}

	return nil
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
	// Server defaults
	if config.Server.Port == 0 {
		config.Server.Port = 8004
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 60 * time.Second
	}

	// Binance defaults
	if config.Binance.BaseURL == "" {
		config.Binance.BaseURL = "https://api.binance.com"
	}
	if config.Binance.Timeout == 0 {
		config.Binance.Timeout = 30 * time.Second
	}
	if config.Binance.MaxRetries == 0 {
		config.Binance.MaxRetries = 3
	}

	// Bitget defaults
	if config.Bitget.BaseURL == "" {
		config.Bitget.BaseURL = "https://api.bitget.com"
	}
	if config.Bitget.Timeout == 0 {
		config.Bitget.Timeout = 30 * time.Second
	}
	if config.Bitget.MaxRetries == 0 {
		config.Bitget.MaxRetries = 3
	}

	// Mobile Money defaults
	if config.MobileMoney.MTN.Timeout == 0 {
		config.MobileMoney.MTN.Timeout = 30 * time.Second
	}
	if config.MobileMoney.MTN.MaxRetries == 0 {
		config.MobileMoney.MTN.MaxRetries = 3
	}
	if config.MobileMoney.Orange.Timeout == 0 {
		config.MobileMoney.Orange.Timeout = 30 * time.Second
	}
	if config.MobileMoney.Orange.MaxRetries == 0 {
		config.MobileMoney.Orange.MaxRetries = 3
	}

	// Resilience defaults
	if config.Resilience.CircuitBreaker.FailureThreshold == 0 {
		config.Resilience.CircuitBreaker.FailureThreshold = 5
	}
	if config.Resilience.CircuitBreaker.RecoveryTimeout == 0 {
		config.Resilience.CircuitBreaker.RecoveryTimeout = 30 * time.Second
	}
	if config.Resilience.CircuitBreaker.Timeout == 0 {
		config.Resilience.CircuitBreaker.Timeout = 30 * time.Second
	}
	if config.Resilience.Retry.MaxRetries == 0 {
		config.Resilience.Retry.MaxRetries = 3
	}
	if config.Resilience.Retry.InitialDelay == 0 {
		config.Resilience.Retry.InitialDelay = 100 * time.Millisecond
	}
	if config.Resilience.Retry.MaxDelay == 0 {
		config.Resilience.Retry.MaxDelay = 5 * time.Second
	}
	if config.Resilience.Timeout.Default == 0 {
		config.Resilience.Timeout.Default = 30 * time.Second
	}
	if config.Resilience.Timeout.API == 0 {
		config.Resilience.Timeout.API = 30 * time.Second
	}

	// Logging defaults
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server configuration
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Validate Binance configuration
	if err := c.Binance.Validate(); err != nil {
		return fmt.Errorf("invalid Binance config: %w", err)
	}

	// Validate Bitget configuration
	if err := c.Bitget.Validate(); err != nil {
		return fmt.Errorf("invalid Bitget config: %w", err)
	}

	// Validate MTN Mobile Money configuration
	if err := c.MobileMoney.MTN.Validate(); err != nil {
		return fmt.Errorf("invalid MTN config: %w", err)
	}

	// Validate Orange Mobile Money configuration
	if err := c.MobileMoney.Orange.Validate(); err != nil {
		return fmt.Errorf("invalid Orange config: %w", err)
	}

	// Validate resilience configuration
	if c.Resilience.CircuitBreaker.FailureThreshold <= 0 {
		return fmt.Errorf("circuit breaker failure threshold must be positive")
	}
	if c.Resilience.CircuitBreaker.RecoveryTimeout <= 0 {
		return fmt.Errorf("circuit breaker recovery timeout must be positive")
	}
	if c.Resilience.Retry.MaxRetries < 0 {
		return fmt.Errorf("retry max retries cannot be negative")
	}
	if c.Resilience.Retry.InitialDelay <= 0 {
		return fmt.Errorf("retry initial delay must be positive")
	}
	if c.Resilience.Retry.MaxDelay <= 0 {
		return fmt.Errorf("retry max delay must be positive")
	}
	if c.Resilience.Timeout.Default <= 0 {
		return fmt.Errorf("default timeout must be positive")
	}
	if c.Resilience.Timeout.API <= 0 {
		return fmt.Errorf("API timeout must be positive")
	}

	// Validate logging configuration
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validFormats[c.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", c.Logging.Format)
	}

	return nil
}

// GetBinanceConfig returns the Binance configuration
func (c *Config) GetBinanceConfig() *models.BinanceConfig {
	return &c.Binance
}

// GetBitgetConfig returns the Bitget configuration
func (c *Config) GetBitgetConfig() *models.BitgetConfig {
	return &c.Bitget
}

// GetMTNConfig returns the MTN Mobile Money configuration
func (c *Config) GetMTNConfig() *models.MobileMoneyConfig {
	return &c.MobileMoney.MTN
}

// GetOrangeConfig returns the Orange Mobile Money configuration
func (c *Config) GetOrangeConfig() *models.MobileMoneyConfig {
	return &c.MobileMoney.Orange
}
