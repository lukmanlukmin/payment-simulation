// Package transaction ...
package transaction

import (
	"context"
	model "payment-simulation/model/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lukmanlukmin/go-lib/database"
)

// Create ...
func (r *Repository) Create(ctx context.Context, trx *model.Transaction) error {
	var db database.SQLQueryExec = r.db
	if tx := database.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	query, args, _ := sq.Insert("transactions").
		Columns("merchant_id", "amount", "direction", "status", "beneficiary_account", "beneficiary_name", "bank_code", "note", "version").
		Values(trx.MerchantID, trx.Amount, trx.Direction, trx.Status, trx.BeneficiaryAccount, trx.BeneficiaryName, trx.BankCode, trx.Note, trx.Version).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := db.QueryRowx(query, args...)
	return row.Scan(&trx.ID, &trx.CreatedAt, &trx.UpdatedAt)
}
