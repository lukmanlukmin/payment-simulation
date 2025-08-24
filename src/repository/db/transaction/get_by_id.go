// Package transaction ...
package transaction

import (
	"context"
	model "payment-simulation/model/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
)

// GetByID ...
func (r *Repository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	tx := &model.Transaction{}
	query, args, _ := sq.Select("id", "merchant_id", "amount", "direction", "status", "beneficiary_account", "beneficiary_name", "bank_code", "note", "version", "created_at", "updated_at").
		From("transactions").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err := db.GetContext(ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
