// Package transactionlog ...
package transactionlog

import (
	"context"
	model "payment-simulation/model/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
)

// Create ...
func (r *Repository) Create(ctx context.Context, log *model.TransactionLog) error {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	query, args, _ := sq.Insert("transaction_logs").
		Columns("transaction_id", "old_status", "new_status", "note").
		Values(log.TransactionID, log.OldStatus, log.NewStatus, log.Note).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := db.ExecContext(ctx, query, args...)
	return err
}
