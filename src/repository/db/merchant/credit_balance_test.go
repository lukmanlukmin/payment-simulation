package merchant

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"payment-simulation/constant"
)

func TestRepository_CreditBalance(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db := sqlx.NewDb(dbMock, "sqlmock")

	repo := &Repository{db: db} // asumsi constructor simple

	ctx := context.Background()
	id := int64(1)
	amount := decimal.NewFromFloat(100)

	testCases := []struct {
		name         string
		expectedRows *sqlmock.Rows
		expectedErr  error
		expectedBal  decimal.Decimal
		expectedVer  int64
	}{
		{
			name:         "success",
			expectedRows: sqlmock.NewRows([]string{"balance", "version"}).AddRow(900.0, 2),
			expectedErr:  nil,
			expectedBal:  decimal.NewFromFloat(900.0),
			expectedVer:  2,
		},
		{
			name:         "insufficient balance",
			expectedRows: nil, // no rows returned
			expectedErr:  constant.ErrTrxInsufficientBalance,
		},
		{
			name:         "db error",
			expectedRows: sqlmock.NewRows([]string{"balance", "version"}),
			expectedErr:  sql.ErrConnDone, // contoh error lain
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := regexp.QuoteMeta(`UPDATE merchants SET balance = balance + $1, version = version + 1, updated_at = NOW() WHERE id = $2 AND version = $3 AND balance >= $4 RETURNING balance, version`)

			if tc.name == "db error" {
				mock.ExpectQuery(query).WithArgs(amount, id, int64(1), amount).WillReturnError(tc.expectedErr)
			} else if tc.name == "insufficient balance" {
				mock.ExpectQuery(query).WithArgs(amount, id, int64(1), amount).WillReturnError(sql.ErrNoRows)
			} else {
				mock.ExpectQuery(query).WithArgs(amount, id, int64(1), amount).WillReturnRows(tc.expectedRows)
			}

			newBal, newVer, err := repo.CreditBalance(ctx, id, amount, 1)
			if tc.expectedErr != nil {
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.True(t, newBal.Equal(tc.expectedBal))
				require.Equal(t, tc.expectedVer, newVer)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
