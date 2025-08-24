// Package event ...
package event

import (
	"context"
	"encoding/json"
	"fmt"
	model "payment-simulation/model/event"
)

// ProcessTransaction ...
func (h *Handler) ProcessTransaction(ctx context.Context, msg []byte) error {
	raw := &model.MessageFormat{}
	err := json.Unmarshal(msg, raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal raw message: %w", err)
	}

	data, _ := json.Marshal(raw.Data)
	updateData := &model.Transaction{}
	err = json.Unmarshal(data, updateData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data message: %w", err)
	}

	return h.Service.TransactionService.ProcessTransaction(ctx, updateData.TransactionID)
}
