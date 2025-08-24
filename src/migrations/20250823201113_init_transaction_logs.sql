-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_logs (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id),
    old_status TEXT CHECK (old_status IN ('PENDING','PROCESSING','SUCCESS','FAILED')),
    new_status TEXT NOT NULL CHECK (new_status IN ('PENDING','PROCESSING','SUCCESS','FAILED')),
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_logs;
-- +goose StatementEnd
