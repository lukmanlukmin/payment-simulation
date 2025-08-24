// Package db ...
package db

import (
	"time"

	"github.com/shopspring/decimal"
)

// Merchant ...
type Merchant struct {
	ID        int64           `db:"id"`
	Name      string          `db:"name"`
	Balance   decimal.Decimal `db:"balance"`
	Version   int64           `db:"version"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
