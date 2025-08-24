// Package cmd ...
package cmd

import (
	"context"
	"payment-simulation/config"
	"payment-simulation/server"
)

// StartHTTP ...
func StartHTTP(_ context.Context, cfg *config.Config) {
	api := server.NewHTTPApi(cfg)
	api.Run()
}
