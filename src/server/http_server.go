// Package server ...
package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"payment-simulation/bootstrap"
	"payment-simulation/config"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HTTPApi ...
type HTTPApi struct {
	cfg *config.Config
}

// NewHTTPApi ...
func NewHTTPApi(cfg *config.Config) *HTTPApi {
	return &HTTPApi{
		cfg: cfg,
	}
}

// Run ...
func (h *HTTPApi) Run() {
	app := fiber.New(fiber.Config{})

	h.HTTPRouter(app, bootstrap.NewBootstrap(h.cfg))
	go func() {
		if err := app.Listen(h.cfg.Server.HTTPPort); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	shutdownGracefully(app)
}

func shutdownGracefully(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("\nShutting down server...")

	// Create a timeout context for cleanup (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown Fiber
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}
