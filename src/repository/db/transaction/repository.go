// Package transaction ...
package transaction

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

import (
	"context"
	model "payment-simulation/model/db"

	"github.com/jmoiron/sqlx"
)

// IRepository ...
type IRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	UpdateStatus(ctx context.Context, tx *model.Transaction) error
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
