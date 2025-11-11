-- +goose Up
-- +goose StatementBegin
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table for authentication component
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_type VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    activated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Create indexes for users table
CREATE INDEX idx_users_status ON users (status);

CREATE INDEX idx_users_user_type ON users (user_type);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- Create user_identifiers table
CREATE TABLE user_identifiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    scheme VARCHAR(50) NOT NULL,
    identifier VARCHAR(255) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMPTZ,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes for user_identifiers table
CREATE INDEX idx_user_identifiers_user_id ON user_identifiers (user_id);

CREATE INDEX idx_user_identifiers_scheme ON user_identifiers (scheme);

CREATE UNIQUE INDEX idx_user_identifiers_scheme_identifier ON user_identifiers (scheme, identifier);

CREATE INDEX idx_user_identifiers_deleted_at ON user_identifiers (deleted_at);

-- Create user_credentials table
CREATE TABLE user_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    scheme VARCHAR(50) NOT NULL,
    credential TEXT NOT NULL,
    expires_at TIMESTAMPTZ,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes for user_credentials table
CREATE INDEX idx_user_credentials_user_id ON user_credentials (user_id);

CREATE INDEX idx_user_credentials_scheme ON user_credentials (scheme);

CREATE UNIQUE INDEX idx_user_credentials_user_id_scheme ON user_credentials (user_id, scheme);

CREATE INDEX idx_user_credentials_deleted_at ON user_credentials (deleted_at);

-- Create user_profiles table
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL UNIQUE,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    email VARCHAR(255),
    birthday DATE,
    locale VARCHAR(10) DEFAULT 'en_US',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes for user_profiles table
CREATE UNIQUE INDEX idx_user_profiles_user_id ON user_profiles (user_id);

CREATE INDEX idx_user_profiles_email ON user_profiles (email);

CREATE INDEX idx_user_profiles_deleted_at ON user_profiles (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS user_profiles;

DROP TABLE IF EXISTS user_credentials;

DROP TABLE IF EXISTS user_identifiers;

DROP TABLE IF EXISTS users;

-- Drop extension (optional, as it might be used by other parts of the system)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd