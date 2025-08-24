// Package transaction ...
package transaction

import (
	"context"
	model "payment-simulation/model/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
)

// UpdateStatus ...
func (r *Repository) UpdateStatus(ctx context.Context, trx *model.Transaction) error {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	query, args, _ := sq.Update("transactions").
		Set("status", trx.Status).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": trx.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := db.ExecContext(ctx, query, args...)
	return err
}
