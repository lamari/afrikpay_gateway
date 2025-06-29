package workflows

import (
	"time"

	"github.com/afrikpay/gateway/internal/config"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Global configuration for workflows
var globalConfig *config.Config

// SetConfig sets the global configuration for workflows
func SetConfig(cfg *config.Config) {
	globalConfig = cfg
}

// GetConfig returns the global configuration
func GetConfig() *config.Config {
	return globalConfig
}

func defaultActivityOptions() workflow.ActivityOptions {
	// TODO: Ajouter des options d'activité par défaut depuis les fichier de configuration
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        0,
		NonRetryableErrorTypes: []string{},
	}
	return workflow.ActivityOptions{
		StartToCloseTimeout:    10 * time.Second,
		ScheduleToCloseTimeout: 100 * time.Second,
		RetryPolicy:            retryPolicy,
	}
}
