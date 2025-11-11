-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id VARCHAR(255),
    username VARCHAR(255),
    action VARCHAR(50),
    success BOOLEAN,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);

CREATE INDEX idx_audit_logs_username ON audit_logs (username);

CREATE INDEX idx_audit_logs_action ON audit_logs (action);

CREATE INDEX idx_audit_logs_success ON audit_logs (success);

CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_audit_logs_created_at;

DROP INDEX IF EXISTS idx_audit_logs_success;

DROP INDEX IF EXISTS idx_audit_logs_action;

DROP INDEX IF EXISTS idx_audit_logs_username;

DROP INDEX IF EXISTS idx_audit_logs_user_id;

DROP TABLE IF EXISTS audit_logs;
-- +goose StatementEnd