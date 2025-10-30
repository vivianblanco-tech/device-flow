-- Create magic_links table
CREATE TABLE IF NOT EXISTS magic_links (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    shipment_id BIGINT REFERENCES shipments(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_magic_links_user_id ON magic_links(user_id);
CREATE INDEX idx_magic_links_token ON magic_links(token);
CREATE INDEX idx_magic_links_expires_at ON magic_links(expires_at);
CREATE INDEX idx_magic_links_shipment_id ON magic_links(shipment_id);
CREATE INDEX idx_magic_links_used_at ON magic_links(used_at);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Comment on tables and columns
COMMENT ON TABLE magic_links IS 'One-time login links sent via email';
COMMENT ON COLUMN magic_links.token IS 'Unique token for the magic link (URL-safe)';
COMMENT ON COLUMN magic_links.expires_at IS 'When the link expires (typically 24-48 hours)';
COMMENT ON COLUMN magic_links.used_at IS 'When the link was used (null if unused)';
COMMENT ON COLUMN magic_links.shipment_id IS 'Optional shipment context for the link';

COMMENT ON TABLE sessions IS 'User sessions for authenticated users';
COMMENT ON COLUMN sessions.token IS 'Unique session token';
COMMENT ON COLUMN sessions.expires_at IS 'When the session expires';

