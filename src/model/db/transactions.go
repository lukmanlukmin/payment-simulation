// Package db ...
package db

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction ...
type Transaction struct {
	ID                 int64           `db:"id"`
	MerchantID         int64           `db:"merchant_id"`
	Amount             decimal.Decimal `db:"amount"`
	Status             string          `db:"status"`    // PENDING, PROCESSING, SUCCESS, FAILED
	Direction          string          `db:"direction"` // DEBIT, CREDIT
	BeneficiaryAccount string          `db:"beneficiary_account"`
	BeneficiaryName    string          `db:"beneficiary_name"`
	BankCode           string          `db:"bank_code"`
	Note               string          `db:"note"`
	Version            int64           `db:"version"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}
