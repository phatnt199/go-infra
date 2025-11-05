-- +goose Up
-- +goose StatementBegin
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    status INT NOT NULL DEFAULT 100,
    user_type VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    activated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    parent_id UUID,
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    CONSTRAINT fk_users_parent FOREIGN KEY (parent_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE INDEX idx_users_status ON users (status)
WHERE
    deleted_at IS NULL;

CREATE INDEX idx_users_parent_id ON users (parent_id)
WHERE
    deleted_at IS NULL;

-- Create user_identifiers table
CREATE TABLE IF NOT EXISTS user_identifiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    scheme VARCHAR(50) NOT NULL DEFAULT 'username',
    identifier VARCHAR(255) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT TRUE,
    details JSONB DEFAULT '{}',
    CONSTRAINT fk_user_identifiers_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_identifiers_deleted_at ON user_identifiers (deleted_at);

CREATE INDEX idx_user_identifiers_user_id ON user_identifiers (user_id)
WHERE
    deleted_at IS NULL;

CREATE UNIQUE INDEX idx_user_identifiers_scheme_identifier ON user_identifiers (scheme, identifier)
WHERE
    deleted_at IS NULL;

-- Create user_credentials table
CREATE TABLE IF NOT EXISTS user_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    scheme VARCHAR(50) NOT NULL DEFAULT 'basic',
    credential VARCHAR(255) NOT NULL,
    details JSONB DEFAULT '{}',
    CONSTRAINT fk_user_credentials_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_credentials_deleted_at ON user_credentials (deleted_at);

CREATE INDEX idx_user_credentials_user_id ON user_credentials (user_id)
WHERE
    deleted_at IS NULL;

CREATE INDEX idx_user_credentials_scheme ON user_credentials (user_id, scheme)
WHERE
    deleted_at IS NULL;

-- Create user_profiles table
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    birthday DATE,
    locale VARCHAR(10) DEFAULT 'en_US',
    details JSONB DEFAULT '{}',
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_profiles_deleted_at ON user_profiles (deleted_at);

CREATE UNIQUE INDEX idx_user_profiles_user_id ON user_profiles (user_id)
WHERE
    deleted_at IS NULL;

-- Create trigger function to update modified_at
CREATE OR REPLACE FUNCTION update_modified_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for all tables
CREATE TRIGGER update_users_modified_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_user_identifiers_modified_at BEFORE UPDATE ON user_identifiers
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_user_credentials_modified_at BEFORE UPDATE ON user_credentials
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_user_profiles_modified_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop triggers
DROP TRIGGER IF EXISTS update_user_profiles_modified_at ON user_profiles;

DROP TRIGGER IF EXISTS update_user_credentials_modified_at ON user_credentials;

DROP TRIGGER IF EXISTS update_user_identifiers_modified_at ON user_identifiers;

DROP TRIGGER IF EXISTS update_users_modified_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_modified_at_column ();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_profiles_user_id;

DROP INDEX IF EXISTS idx_user_profiles_deleted_at;

DROP INDEX IF EXISTS idx_user_credentials_scheme;

DROP INDEX IF EXISTS idx_user_credentials_user_id;

DROP INDEX IF EXISTS idx_user_credentials_deleted_at;

DROP INDEX IF EXISTS idx_user_identifiers_scheme_identifier;

DROP INDEX IF EXISTS idx_user_identifiers_user_id;

DROP INDEX IF EXISTS idx_user_identifiers_deleted_at;

DROP INDEX IF EXISTS idx_users_parent_id;

DROP INDEX IF EXISTS idx_users_status;

DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop tables
DROP TABLE IF EXISTS user_profiles;

DROP TABLE IF EXISTS user_credentials;

DROP TABLE IF EXISTS user_identifiers;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd