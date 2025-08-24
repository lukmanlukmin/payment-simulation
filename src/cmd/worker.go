// Package cmd ...
package cmd

import (
	"context"
	"payment-simulation/config"
	"payment-simulation/server"
)

// StartWorker ...
func StartWorker(_ context.Context, cfg *config.Config) {
	worker := server.NewEventServer(cfg)
	worker.Run()
}
