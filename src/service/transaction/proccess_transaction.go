// Package transaction ...
package transaction

import (
	"context"
	"math/rand"
	"payment-simulation/constant"
	"payment-simulation/model/db"
	"time"

	"github.com/lukmanlukmin/go-lib/database"
)

// ProcessTransaction ...
func (s *Service) ProcessTransaction(ctx context.Context, transactionID int64) error {
	trx, err := s.TransactionRepository.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if trx.Status != constant.TrxStatusPending {
		return constant.ErrTrxAlreadyProcessed
	}

	oldStatus := trx.Status
	trx.Status = constant.TrxStatusProcessing
	err = database.BeginTransaction(ctx, s.Store.GetMaster(), func(ctx context.Context) error {
		if err := s.TransactionRepository.UpdateStatus(ctx, trx); err != nil {
			return err
		}
		if err := s.TransactionLogRepository.Create(ctx, &db.TransactionLog{
			TransactionID: trx.ID,
			OldStatus:     oldStatus,
			NewStatus:     trx.Status,
			Note:          "Processing transaction",
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// simulate provider
	time.Sleep(time.Duration(300+rand.Intn(1200)) * time.Millisecond)
	success := rand.Intn(100) < 98 // 98% chance success

	if success {
		oldStatus = trx.Status
		trx.Status = constant.TrxStatusSuccess
		err = database.BeginTransaction(ctx, s.Store.GetMaster(), func(ctx context.Context) error {
			if err := s.TransactionRepository.UpdateStatus(ctx, trx); err != nil {
				return err
			}
			if err := s.TransactionLogRepository.Create(ctx, &db.TransactionLog{
				TransactionID: trx.ID,
				OldStatus:     oldStatus,
				NewStatus:     trx.Status,
				Note:          "Finishing transaction",
			}); err != nil {
				return err
			}
			return nil
		})
		return err
	}

	const maxRetries = 5
	backoff := 10 * time.Millisecond
	for attempt := 0; attempt < maxRetries; attempt++ {
		merchant, err := s.MerchantRepository.GetByID(ctx, trx.MerchantID)
		if err != nil {
			return err
		}
		err = database.BeginTransaction(ctx, s.Store.GetMaster(), func(ctx context.Context) error {
			_, _, err = s.MerchantRepository.CreditBalance(ctx, trx.MerchantID, trx.Amount, merchant.Version)
			if err != nil {
				return err
			}
			oldStatus = trx.Status
			trx.Status = constant.TrxStatusFailed
			if err := s.TransactionRepository.UpdateStatus(ctx, trx); err != nil {
				return err
			}
			if err := s.TransactionLogRepository.Create(ctx, &db.TransactionLog{
				TransactionID: trx.ID,
				OldStatus:     oldStatus,
				NewStatus:     trx.Status,
				Note:          "Transaction 3rd Party Error",
			}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			time.Sleep(backoff + time.Duration(rand.Intn(5))*time.Millisecond)
			backoff *= 2
			continue
		}
		return err
	}
	return nil
}
