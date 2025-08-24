// Package merchant ...
package merchant

import (
	"context"
	"database/sql"
	"errors"
	model "payment-simulation/model/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
)

// GetByID ...
func (r *Repository) GetByID(ctx context.Context, id int64) (*model.Merchant, error) {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	query, args, _ := sq.Select("id", "name", "balance", "version", "created_at", "updated_at").
		From("merchants").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var m model.Merchant
	if err := db.GetContext(ctx, &m, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
