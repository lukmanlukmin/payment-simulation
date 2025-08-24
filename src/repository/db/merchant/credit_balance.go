// Package merchant ...
package merchant

import (
	"context"
	"database/sql"
	"payment-simulation/constant"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
	"github.com/shopspring/decimal"
)

// CreditBalance ...
func (r *Repository) CreditBalance(ctx context.Context, id int64, amount decimal.Decimal, expectedVersion int64) (newBalance decimal.Decimal, newVersion int64, err error) {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	// Build query dengan RETURNING
	query, args, _ := sq.Update("merchants").
		Set("balance", sq.Expr("balance + ?", amount)).
		Set("version", sq.Expr("version + 1")).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id, "version": expectedVersion}).
		Where(sq.GtOrEq{"balance": amount}).
		Suffix("RETURNING balance, version").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	// Eksekusi query
	row := db.QueryRowx(query, args...)
	if err = row.Scan(&newBalance, &newVersion); err != nil {
		if err == sql.ErrNoRows {
			err = constant.ErrTrxInsufficientBalance
		}
		return
	}

	return
}
