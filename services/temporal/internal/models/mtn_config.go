package models

import (
	"time"
)

// MTNConfig holds configuration for MTN Mobile Money API client
type MTNConfig struct {
	BaseURL      string        `yaml:"base_url"`
	PrimaryKey   string        `yaml:"primary_key"`
	SecondaryKey string        `yaml:"secondary_key"`
	Timeout      time.Duration `yaml:"timeout"`
	MaxRetries   int           `yaml:"rate_limit"`
}

// Validate validates the MTN configuration
func (c *MTNConfig) Validate() error {
	if c.BaseURL == "" {
		return NewClientError("INVALID_CONFIG", "base URL is required", false)
	}
	if c.PrimaryKey == "" {
		return NewClientError("INVALID_CONFIG", "primary key is required", false)
	}
	if c.SecondaryKey == "" {
		return NewClientError("INVALID_CONFIG", "secondary key is required", false)
	}
	if c.Timeout <= 0 {
		return NewClientError("INVALID_CONFIG", "timeout must be positive", false)
	}
	if c.MaxRetries < 0 {
		return NewClientError("INVALID_CONFIG", "max retries cannot be negative", false)
	}
	return nil
}
