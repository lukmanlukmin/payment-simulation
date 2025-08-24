// Package repository ...
package repository

import (
	"payment-simulation/repository/db/merchant"
	"payment-simulation/repository/db/transaction"
	transactionlog "payment-simulation/repository/db/transaction_log"
	"payment-simulation/repository/event"

	connDB "github.com/lukmanlukmin/go-lib/database/connection"
	"github.com/lukmanlukmin/go-lib/kafka"
)

// Repository ...
type Repository struct {
	Store                    *connDB.Store
	MerchantRepository       merchant.IRepository
	TransactionRepository    transaction.IRepository
	TransactionLogRepository transactionlog.IRepository
	KafkaProducer            event.IRepository
}

// LoadRepository ...
func LoadRepository(storeDB *connDB.Store, kafkaProducer kafka.Producer) *Repository {

	return &Repository{
		Store:                    storeDB,
		MerchantRepository:       merchant.NewRepository(storeDB.GetMaster()),
		TransactionRepository:    transaction.NewRepository(storeDB.GetMaster()),
		TransactionLogRepository: transactionlog.NewRepository(storeDB.GetMaster()),
		KafkaProducer:            event.NewRepository(kafkaProducer),
	}
}
