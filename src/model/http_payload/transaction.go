// Package httppayload ...
package httppayload

import (
	"github.com/shopspring/decimal"
)

// TransferRequest ...
type TransferRequest struct {
	MerchantID         int64   `swaggerignore:"true"`
	Amount             float64 `json:"amount" validate:"required,gt=0"` // must be > 0
	BeneficiaryName    string  `json:"beneficiary_name" validate:"required,min=1,max=100"`
	BeneficiaryAccount string  `json:"beneficiary_account" validate:"required,min=1,max=50"`
	BankCode           string  `json:"bank_code" validate:"required,min=1,max=20"`
	Note               string  `json:"note,omitempty" validate:"max=255"` // optional
}

// TransferResponse ...
type TransferResponse struct {
	TransactionID      int64           `json:"transaction_id"`
	Amount             decimal.Decimal `json:"amount"`
	Status             string          `json:"status"`
	BeneficiaryName    string          `json:"beneficiary_name"`
	BeneficiaryAccount string          `json:"beneficiary_account"`
	BankCode           string          `json:"bank_code"`
	Note               string          `json:"note"`
}
