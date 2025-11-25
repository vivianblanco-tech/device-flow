-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    role user_role NOT NULL,
    google_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- Add constraint to ensure either password_hash or google_id is provided
ALTER TABLE users ADD CONSTRAINT chk_users_auth_method 
    CHECK (
        (password_hash IS NOT NULL AND password_hash != '') OR 
        (google_id IS NOT NULL AND google_id != '')
    );

-- Comment on table and columns
COMMENT ON TABLE users IS 'User accounts for Align';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password (for password auth)';
COMMENT ON COLUMN users.role IS 'User role: logistics, client, warehouse, or project_manager';
COMMENT ON COLUMN users.google_id IS 'Google OAuth ID (for Google auth)';

