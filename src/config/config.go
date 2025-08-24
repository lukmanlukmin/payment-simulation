// Package config ...
package config

import (
	"payment-simulation/constant"

	fileConfig "github.com/lukmanlukmin/go-lib/file"
	kafka "github.com/lukmanlukmin/go-lib/kafka"
	"github.com/lukmanlukmin/go-lib/log"
)

// Config ...
type Config struct {
	Server           ServerConfig   `yaml:"Server"`
	PostgreSQLConfig PostgresConfig `yaml:"PostgresDBConfig"`
	KafkaConfig      kafka.Config   `yaml:"KafkaConfig"`
	Topics           Topics         `yaml:"Topics"`
}

type (
	// ServerConfig Configuration
	ServerConfig struct {
		HTTPPort string `yaml:"HttpPort"`
	}

	// Security Configuration
	Security struct {
		JWTSecret   string `yaml:"JWTSecret"`
		JWTDuration string `yaml:"JWTDuration"`
	}

	// PostgresConfig Configuration
	PostgresConfig struct {
		DSN           string `yaml:"DSN"`
		RetryInterval int    `yaml:"RetryInterval"`
		MaxIdleConn   int    `yaml:"MaxIdleConn"`
		MaxConn       int    `yaml:"MaxConn"`
	}

	// Topics Configuration
	Topics struct {
		TransactionTopic string `yaml:"TransactionTopic"`
	}
)

// ReadModuleConfig ...
func ReadModuleConfig(cfg interface{}, filePath string) error {
	if filePath != "" {
		err := fileConfig.ReadConfig(cfg, filePath)
		if err != nil {
			log.Fatalf("failed to read config. %v", err)
		}
		return nil
	}
	err := fileConfig.ReadConfig(cfg, constant.DefaultConfigFile)
	if err != nil {
		log.Fatalf("failed to read config. %v", err)
	}
	return nil
}
