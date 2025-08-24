// Package db ...
package db

import "time"

// TransactionLog ...
type TransactionLog struct {
	ID            int64     `db:"id"`
	TransactionID int64     `db:"transaction_id"`
	OldStatus     string    `db:"old_status"`
	NewStatus     string    `db:"new_status"`
	Note          string    `db:"note"`
	CreatedAt     time.Time `db:"created_at"`
}
