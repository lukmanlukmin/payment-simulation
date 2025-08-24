-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    amount NUMERIC(18,2) NOT NULL CHECK (amount > 0),
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('DEBIT', 'CREDIT')),
    status TEXT NOT NULL CHECK (status IN ('PENDING','PROCESSING','SUCCESS','FAILED')),
    beneficiary_account VARCHAR(50),      -- rekening tujuan/pengirim
    beneficiary_name VARCHAR(100),        -- nama penerima/pengirim
    bank_code VARCHAR(20),                -- kode bank
    note TEXT,                            -- opsional, catatan tambahan
    version INT NOT NULL DEFAULT 1,       -- buat optimistic locking
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
