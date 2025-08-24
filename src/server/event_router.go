package server

import (
	"context"
	"payment-simulation/bootstrap"
	"payment-simulation/config"

	"payment-simulation/handler/event"

	kafkalib "github.com/lukmanlukmin/go-lib/kafka"
)

// StartConsumers ...
func StartConsumers(ctx context.Context, b *bootstrap.Bootstrap, cfg *config.Config) {

	consumer := kafkalib.NewConsumerGroup(&cfg.KafkaConfig)

	handler := event.NewHandler(b)
	consumer.Subscribe(&kafkalib.ConsumerContext{
		Context: ctx,
		Topics:  []string{cfg.Topics.TransactionTopic},
		GroupID: cfg.KafkaConfig.ClientID,
		Handler: kafkalib.MessageProcessorFunc(func(msg *kafkalib.MessageDecoder) {
			if err := handler.ProcessTransaction(ctx, msg.Body); err != nil {
				msg.Commit(msg)
			}
		}),
	})

}
