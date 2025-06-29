package models

import (
	"time"
)

// OrangeConfig holds configuration for Orange Money API client
type OrangeConfig struct {
	BaseURL       string        `yaml:"base_url"`
	ClientID      string        `yaml:"client_id"`
	ClientSecret  string        `yaml:"client_secret"`
	Authorization string        `yaml:"authorization"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxRetries    int           `yaml:"rate_limit"`
}

// Validate validates the Orange configuration
func (c *OrangeConfig) Validate() error {
	if c.BaseURL == "" {
		return NewClientError("INVALID_CONFIG", "base URL is required", false)
	}
	if c.ClientID == "" {
		return NewClientError("INVALID_CONFIG", "client ID is required", false)
	}
	if c.ClientSecret == "" {
		return NewClientError("INVALID_CONFIG", "client secret is required", false)
	}
	if c.Authorization == "" {
		return NewClientError("INVALID_CONFIG", "authorization is required", false)
	}
	if c.Timeout <= 0 {
		return NewClientError("INVALID_CONFIG", "timeout must be positive", false)
	}
	if c.MaxRetries < 0 {
		return NewClientError("INVALID_CONFIG", "max retries cannot be negative", false)
	}
	return nil
}
