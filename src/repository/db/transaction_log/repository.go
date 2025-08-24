// Package transactionlog ...
package transactionlog

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

import (
	"context"
	model "payment-simulation/model/db"

	"github.com/jmoiron/sqlx"
)

// IRepository ...
type IRepository interface {
	Create(ctx context.Context, log *model.TransactionLog) error
}

// Repository ...
type Repository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
