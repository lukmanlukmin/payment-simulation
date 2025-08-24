// Package event ...
package event

import (
	"context"
	model "payment-simulation/model/event"

	"github.com/lukmanlukmin/go-lib/kafka"
)

// Publish ...
func (r *Repository) Publish(ctx context.Context, topic string, value interface{}) error {
	data, err := model.BuildKafkaPayload(value, topic)
	if err != nil {
		return err
	}
	return r.Producer.Publish(ctx, &kafka.MessageContext{
		Topic: topic,
		Value: data,
	})
}
