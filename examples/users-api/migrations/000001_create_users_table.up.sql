-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users (email)
WHERE
    deleted_at IS NULL;

-- Insert sample data
INSERT INTO
    users (name, email)
VALUES (
        'Alice Johnson',
        'alice@example.com'
    ),
    (
        'Bob Smith',
        'bob@example.com'
    ),
    (
        'Charlie Brown',
        'charlie@example.com'
    ),
    (
        'Diana Prince',
        'diana@example.com'
    ),
    (
        'Eve Adams',
        'eve@example.com'
    );