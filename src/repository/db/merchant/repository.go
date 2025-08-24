// Package merchant ...
package merchant

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

import (
	"context"
	model "payment-simulation/model/db"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// IRepository ...
type IRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Merchant, error)
	DeductBalance(ctx context.Context, id int64, amount decimal.Decimal, expectedVersion int64) (newBalance decimal.Decimal, newVersion int64, err error)
	CreditBalance(ctx context.Context, id int64, amount decimal.Decimal, expectedVersion int64) (newBalance decimal.Decimal, newVersion int64, err error)
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
