package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"payment-simulation/bootstrap/repository"
	"payment-simulation/config"
	"payment-simulation/constant"
	model "payment-simulation/model/db"
	eventModel "payment-simulation/model/event"
	payload "payment-simulation/model/http_payload"
	mockMerchant "payment-simulation/repository/db/merchant/mocks"
	mockTransaction "payment-simulation/repository/db/transaction/mocks"
	mockKafka "payment-simulation/repository/event/mocks"

	connDB "github.com/lukmanlukmin/go-lib/database/connection"
)

func TestService_SubmitTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// sqlmock + sqlx
	dbStd, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	dbx := sqlx.NewDb(dbStd, "sqlmock")
	store := &connDB.Store{
		Master: dbx,
		Slave:  dbx,
	}

	// Mock dependencies
	mockMerchantRepo := mockMerchant.NewMockIRepository(ctrl)
	mockTransactionRepo := mockTransaction.NewMockIRepository(ctrl)
	mockKafka := mockKafka.NewMockIRepository(ctrl)

	tests := []struct {
		name       string
		req        payload.TransferRequest
		setupMocks func()
		wantErr    error
		wantStatus string
	}{
		{
			name: "success transfer",
			req: payload.TransferRequest{
				MerchantID:         1,
				Amount:             100,
				BeneficiaryName:    "Alice",
				BeneficiaryAccount: "123",
				BankCode:           "BCA",
				Note:               "test",
			},
			setupMocks: func() {
				merchant := &model.Merchant{ID: 1, Balance: decimal.NewFromInt(200), Version: 1}

				// sqlmock expect transaction
				mockDB.ExpectBegin()

				// gomock expectations, gunakan AnyTimes() karena retry loop
				mockMerchantRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(merchant, nil).AnyTimes()
				mockMerchantRepo.EXPECT().DeductBalance(gomock.Any(), int64(1), decimal.NewFromFloat(100), int64(1)).
					Return(decimal.NewFromInt(100), int64(2), nil).AnyTimes()

				mockTransactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, trx *model.Transaction) error {
						trx.ID = 10
						return nil
					}).AnyTimes()

				mockKafka.EXPECT().Publish(gomock.Any(), "trx-topic", gomock.AssignableToTypeOf(eventModel.Transaction{})).
					Return(nil).AnyTimes()
				mockDB.ExpectCommit()
			},
			wantErr:    nil,
			wantStatus: constant.TrxStatusPending,
		},
		{
			name: "insufficient balance",
			req: payload.TransferRequest{
				MerchantID:         1,
				Amount:             500,
				BeneficiaryName:    "Bob",
				BeneficiaryAccount: "456",
				BankCode:           "BNI",
			},
			setupMocks: func() {
				merchant := &model.Merchant{ID: 1, Balance: decimal.NewFromInt(100), Version: 1}
				mockMerchantRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(merchant, nil).AnyTimes()
			},
			wantErr:    constant.ErrTrxInsufficientBalance,
			wantStatus: "",
		},
		{
			name: "repo error",
			req: payload.TransferRequest{
				MerchantID: 1,
				Amount:     50,
			},
			setupMocks: func() {
				mockMerchantRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, errors.New("db error")).AnyTimes()
			},
			wantErr:    constant.ErrTrxBusy,
			wantStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			svc := NewService(&repository.Repository{
				Store:                 store,
				MerchantRepository:    mockMerchantRepo,
				TransactionRepository: mockTransactionRepo,
				KafkaProducer:         mockKafka,
			}, &config.Config{
				Topics: config.Topics{
					TransactionTopic: "trx-topic",
				},
			})
			resp, err := svc.SubmitTransfer(ctx, tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantStatus, resp.Status)
				assert.Equal(t, int64(10), resp.TransactionID)
			}

			// pastikan semua sqlmock expectations terpenuhi
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
