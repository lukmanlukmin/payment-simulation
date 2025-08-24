// Package cmd ...
package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"payment-simulation/config"

	"github.com/spf13/cobra"
)

// Start ...
func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	rootCmd := &cobra.Command{}
	serveHTTPCmd := &cobra.Command{
		Use:   "serve-http",
		Short: "Run HTTP Server",
		Run: func(_ *cobra.Command, _ []string) {
			cfg := &config.Config{}
			err := config.ReadModuleConfig(cfg, "") /// set default config // working on it later
			if err == nil {
				StartHTTP(ctx, cfg)
			}
		},
	}
	rootCmd.AddCommand(serveHTTPCmd)

	serveWorkerCmd := &cobra.Command{
		Use:   "serve-worker",
		Short: "Run Worker",
		Run: func(_ *cobra.Command, _ []string) {
			cfg := &config.Config{}
			err := config.ReadModuleConfig(cfg, "") /// set default config // working on it later
			if err == nil {
				StartWorker(ctx, cfg)
			}
		},
	}
	rootCmd.AddCommand(serveWorkerCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
