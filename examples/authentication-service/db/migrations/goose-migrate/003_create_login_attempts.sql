-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    identifier VARCHAR(255) NOT NULL,
    attempts INT DEFAULT 0,
    locked_until TIMESTAMPTZ,
    last_attempt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_login_attempts_identifier ON login_attempts (identifier);

CREATE INDEX idx_login_attempts_locked_until ON login_attempts (locked_until);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_login_attempts_locked_until;

DROP INDEX IF EXISTS idx_login_attempts_identifier;

DROP TABLE IF EXISTS login_attempts;
-- +goose StatementEnd