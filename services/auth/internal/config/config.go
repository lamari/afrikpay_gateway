package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the auth service
type Config struct {
	Server ServerConfig `yaml:"server"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	PrivateKeyPath     string        `yaml:"private_key_path"`
	PublicKeyPath      string        `yaml:"public_key_path"`
	AccessTokenExpiry  time.Duration `yaml:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `yaml:"refresh_token_expiry"`
	Issuer             string        `yaml:"issuer"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("AUTH_PORT", 8001),
			Host:         getEnv("AUTH_HOST", "0.0.0.0"),
			ReadTimeout:  time.Duration(getEnvAsInt("AUTH_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("AUTH_WRITE_TIMEOUT", 10)) * time.Second,
			IdleTimeout:  time.Duration(getEnvAsInt("AUTH_IDLE_TIMEOUT", 60)) * time.Second,
		},
		JWT: JWTConfig{
			PrivateKeyPath:     getEnv("JWT_PRIVATE_KEY_PATH", "../../config/keys/private.pem"),
			PublicKeyPath:      getEnv("JWT_PUBLIC_KEY_PATH", "../../config/keys/public.pem"),
			AccessTokenExpiry:  time.Duration(getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRY", 15)) * time.Minute,
			RefreshTokenExpiry: time.Duration(getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRY", 24*7)) * time.Hour,
			Issuer:             getEnv("JWT_ISSUER", "afrikpay-gateway"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.JWT.PrivateKeyPath == "" {
		return fmt.Errorf("JWT private key path is required")
	}

	if c.JWT.PublicKeyPath == "" {
		return fmt.Errorf("JWT public key path is required")
	}

	if c.JWT.AccessTokenExpiry <= 0 {
		return fmt.Errorf("JWT access token expiry must be positive")
	}

	if c.JWT.RefreshTokenExpiry <= 0 {
		return fmt.Errorf("JWT refresh token expiry must be positive")
	}

	return nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
