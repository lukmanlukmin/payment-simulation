// Package service ...
package service

import (
	"payment-simulation/bootstrap/repository"
	"payment-simulation/config"
	"payment-simulation/service/transaction"

	connDB "github.com/lukmanlukmin/go-lib/database/connection"
)

// Service ...
type Service struct {
	Store              *connDB.Store
	TransactionService transaction.IService
}

// LoadService ...
func LoadService(bs *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Store:              bs.Store,
		TransactionService: transaction.NewService(bs, cfg),
	}
}
