// Package server ...
package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payment-simulation/bootstrap"
	"payment-simulation/config"
)

// EventServer ...
type EventServer struct {
	cfg *config.Config
}

// NewEventServer ...
func NewEventServer(cfg *config.Config) *EventServer {
	return &EventServer{cfg: cfg}
}

// Run ...
func (s *EventServer) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumers
	StartConsumers(ctx, bootstrap.NewBootstrap(s.cfg), s.cfg)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Stopping event consumers...")

	// Delay to allow consumers finish processing
	time.Sleep(2 * time.Second)
	cancel()
	log.Println("Event consumer shutdown completed.")
}
