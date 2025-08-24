package constant

import "errors"

var (
	// ErrTrxInsufficientBalance ...
	ErrTrxInsufficientBalance = errors.New("insufficient balance")
	// ErrTrxConcurrentUpdateDetected ...
	ErrTrxConcurrentUpdateDetected = errors.New("concurrent update detected")
	// ErrTrxBusy ...
	ErrTrxBusy = errors.New("transaction too busy, please retry later")
	// ErrTrxAlreadyProcessed ...
	ErrTrxAlreadyProcessed = errors.New("transaction already processed")
)
