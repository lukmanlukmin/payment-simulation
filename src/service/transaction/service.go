// Package transaction ...
package transaction

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

import (
	"context"
	"payment-simulation/bootstrap/repository"
	"payment-simulation/config"
	payload "payment-simulation/model/http_payload"
)

// Service ...
type Service struct {
	*repository.Repository
	cfg *config.Config
}

// NewService ...
func NewService(bs *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Repository: bs,
		cfg:        cfg,
	}
}

// IService ...
type IService interface {
	// Submit a new transfer request
	// Deducts merchant balance atomically using optimistic locking
	SubmitTransfer(ctx context.Context, req payload.TransferRequest) (*payload.TransferResponse, error)

	// Process transaction asynchronously (to be called by worker)
	// Updates transaction status and logs changes
	ProcessTransaction(ctx context.Context, transactionID int64) error
}
