-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchants (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    balance NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    version INT NOT NULL DEFAULT 1, 
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
INSERT INTO merchants (id, name, balance) VALUES (1, 'Default Merchant', 100000000.00)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merchants;
-- +goose StatementEnd
