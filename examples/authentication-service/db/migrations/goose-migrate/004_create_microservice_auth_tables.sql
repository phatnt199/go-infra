-- +goose Up
-- +goose StatementBegin
-- Create users table for microservice-auth
CREATE TABLE IF NOT EXISTS ms_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_type VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ms_users_status ON ms_users (status);

CREATE INDEX idx_ms_users_user_type ON ms_users (user_type);

-- Create user_identifiers table
CREATE TABLE IF NOT EXISTS ms_user_identifiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    scheme VARCHAR(50) NOT NULL,
    identifier VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES ms_users (id) ON DELETE CASCADE
);

CREATE INDEX idx_ms_user_identifiers_user_id ON ms_user_identifiers (user_id);

CREATE INDEX idx_ms_user_identifiers_scheme ON ms_user_identifiers (scheme);

CREATE UNIQUE INDEX idx_ms_user_identifiers_scheme_identifier ON ms_user_identifiers (scheme, identifier);

-- Create user_credentials table
CREATE TABLE IF NOT EXISTS ms_user_credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    scheme VARCHAR(50) NOT NULL,
    credential TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES ms_users (id) ON DELETE CASCADE
);

CREATE INDEX idx_ms_user_credentials_user_id ON ms_user_credentials (user_id);

CREATE INDEX idx_ms_user_credentials_scheme ON ms_user_credentials (scheme);

-- Create user_profiles table
CREATE TABLE IF NOT EXISTS ms_user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL UNIQUE,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    email VARCHAR(255),
    birthday DATE,
    locale VARCHAR(10) DEFAULT 'en_US',
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES ms_users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_ms_user_profiles_user_id ON ms_user_profiles (user_id);

CREATE INDEX idx_ms_user_profiles_email ON ms_user_profiles (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ms_user_profiles_email;

DROP INDEX IF EXISTS idx_ms_user_profiles_user_id;

DROP TABLE IF EXISTS ms_user_profiles;

DROP INDEX IF EXISTS idx_ms_user_credentials_scheme;

DROP INDEX IF EXISTS idx_ms_user_credentials_user_id;

DROP TABLE IF EXISTS ms_user_credentials;

DROP INDEX IF EXISTS idx_ms_user_identifiers_scheme_identifier;

DROP INDEX IF EXISTS idx_ms_user_identifiers_scheme;

DROP INDEX IF EXISTS idx_ms_user_identifiers_user_id;

DROP TABLE IF EXISTS ms_user_identifiers;

DROP INDEX IF EXISTS idx_ms_users_user_type;

DROP INDEX IF EXISTS idx_ms_users_status;

DROP TABLE IF EXISTS ms_users;
-- +goose StatementEnd