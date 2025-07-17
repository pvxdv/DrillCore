-- +goose Up
CREATE TABLE IF NOT EXISTS debt (
id SERIAL PRIMARY KEY,
user_id BIGINT NOT NULL,
description TEXT,
amount BIGINT NOT NULL CHECK (amount > 0),
return_date TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS debt_user_id_idx ON debt(user_id);

-- +goose Down
DROP TABLE IF EXISTS debt;