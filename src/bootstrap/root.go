// Package bootstrap ...
package bootstrap

import (
	"payment-simulation/bootstrap/repository"
	"payment-simulation/bootstrap/service"
	"payment-simulation/config"

	connDB "github.com/lukmanlukmin/go-lib/database/connection"
	"github.com/lukmanlukmin/go-lib/kafka"
)

// Bootstrap ...
type Bootstrap struct {
	Repository *repository.Repository
	Service    *service.Service
}

// NewBootstrap ...
func NewBootstrap(conf *config.Config) *Bootstrap {
	storeDB := connDB.New(connDB.DBConfig{
		MasterDSN:     conf.PostgreSQLConfig.DSN,
		EnableSlave:   false,
		RetryInterval: conf.PostgreSQLConfig.RetryInterval,
		MaxIdleConn:   conf.PostgreSQLConfig.MaxIdleConn,
		MaxConn:       conf.PostgreSQLConfig.MaxConn,
	}, connDB.DriverPostgres)
	kafkaProducer := kafka.NewProducer(&conf.KafkaConfig)

	repo := repository.LoadRepository(storeDB, kafkaProducer)
	svc := service.LoadService(repo, conf)

	return &Bootstrap{
		Repository: repo,
		Service:    svc,
	}
}
