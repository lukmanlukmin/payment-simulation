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
	mockMerchant "payment-simulation/repository/db/merchant/mocks"
	mockTransaction "payment-simulation/repository/db/transaction/mocks"
	mockTransactionLog "payment-simulation/repository/db/transaction_log/mocks"

	connDB "github.com/lukmanlukmin/go-lib/database/connection"
)

func TestService_ProcessTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	dbStd, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	dbx := sqlx.NewDb(dbStd, "sqlmock")
	store := &connDB.Store{Master: dbx, Slave: dbx}

	mockTransactionRepo := mockTransaction.NewMockIRepository(ctrl)
	mockTransactionLogRepo := mockTransactionLog.NewMockIRepository(ctrl)
	mockMerchantRepo := mockMerchant.NewMockIRepository(ctrl)

	svc := NewService(&repository.Repository{
		Store:                    store,
		TransactionRepository:    mockTransactionRepo,
		TransactionLogRepository: mockTransactionLogRepo,
		MerchantRepository:       mockMerchantRepo,
	}, &config.Config{})

	trxPending := &model.Transaction{ID: 1, MerchantID: 10, Amount: decimal.NewFromInt(100), Status: constant.TrxStatusPending}
	trxProcessed := &model.Transaction{ID: 2, MerchantID: 11, Amount: decimal.NewFromInt(50), Status: constant.TrxStatusSuccess}

	tests := []struct {
		name       string
		trxID      int64
		setupMocks func()
		wantErr    error
	}{
		{
			name:  "transaction not found",
			trxID: 99,
			setupMocks: func() {
				mockTransactionRepo.EXPECT().GetByID(gomock.Any(), int64(99)).Return(nil, errors.New("not found"))
			},
			wantErr: errors.New("not found"),
		},
		{
			name:  "already processed",
			trxID: 2,
			setupMocks: func() {
				mockTransactionRepo.EXPECT().GetByID(gomock.Any(), int64(2)).Return(trxProcessed, nil)
			},
			wantErr: constant.ErrTrxAlreadyProcessed,
		},
		{
			name:  "success transaction",
			trxID: 1,
			setupMocks: func() {
				// sqlmock transaction
				mockDB.ExpectBegin()
				mockDB.ExpectCommit()
				mockDB.ExpectBegin()
				mockDB.ExpectCommit()

				mockTransactionRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(trxPending, nil)
				mockTransactionRepo.EXPECT().UpdateStatus(gomock.Any(), trxPending).Return(nil).AnyTimes()
				mockTransactionLogRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.ProcessTransaction(ctx, tt.trxID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
