package transaction

import (
	"context"
	"math/rand"
	"payment-simulation/constant"
	model "payment-simulation/model/db"
	eventModel "payment-simulation/model/event"
	payload "payment-simulation/model/http_payload"
	"time"

	"github.com/lukmanlukmin/go-lib/database"
	"github.com/shopspring/decimal"
)

// SubmitTransfer ...
func (s *Service) SubmitTransfer(ctx context.Context, req payload.TransferRequest) (*payload.TransferResponse, error) {
	const maxRetries = 5
	backoff := 10 * time.Millisecond

	trxID := int64(0)
	for attempt := 0; attempt < maxRetries; attempt++ {
		merchant, err := s.MerchantRepository.GetByID(ctx, req.MerchantID)
		if err != nil {
			return nil, err
		}
		if merchant.Balance.LessThan(decimal.NewFromFloat(req.Amount)) {
			return nil, constant.ErrTrxInsufficientBalance
		}

		err = database.BeginTransaction(ctx, s.Store.GetMaster(), func(ctx context.Context) error {
			_, newVersion, err := s.MerchantRepository.DeductBalance(ctx, req.MerchantID, decimal.NewFromFloat(req.Amount), merchant.Version)
			if err != nil {
				return err
			}

			trx := &model.Transaction{
				MerchantID:         req.MerchantID,
				Amount:             decimal.NewFromFloat(req.Amount),
				Direction:          constant.TrxDirectionDebit,
				Status:             constant.TrxStatusPending,
				BeneficiaryName:    req.BeneficiaryName,
				BeneficiaryAccount: req.BeneficiaryAccount,
				BankCode:           req.BankCode,
				Note:               req.Note,
				Version:            newVersion,
			}
			if err := s.TransactionRepository.Create(ctx, trx); err != nil {
				return err
			}
			trxID = trx.ID
			err = s.KafkaProducer.Publish(ctx, s.cfg.Topics.TransactionTopic, eventModel.Transaction{
				TransactionID: trx.ID,
			})
			return err
		})
		if err != nil {
			time.Sleep(backoff + time.Duration(rand.Intn(5))*time.Millisecond)
			backoff *= 2
			continue
		}
		result := &payload.TransferResponse{
			TransactionID:      trxID,
			Amount:             decimal.NewFromFloat(req.Amount),
			Status:             constant.TrxStatusPending,
			BeneficiaryName:    req.BeneficiaryName,
			BeneficiaryAccount: req.BeneficiaryAccount,
			BankCode:           req.BankCode,
			Note:               req.Note,
		}
		return result, nil
	}

	return nil, constant.ErrTrxBusy
}
