package workflows

import (
	"time"

	"github.com/afrikpay/gateway/internal/config"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var configInstance *config.Config

func SetConfig(cfg *config.Config) {
	configInstance = cfg
}

func GetConfig() *config.Config {
	return configInstance
}

// options d'activité par défaut
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

// WorkflowHandler gère les requêtes de workflow
