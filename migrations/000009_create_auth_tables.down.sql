-- Drop indexes for sessions
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_sessions_user_id;

-- Drop indexes for magic_links
DROP INDEX IF EXISTS idx_magic_links_used_at;
DROP INDEX IF EXISTS idx_magic_links_shipment_id;
DROP INDEX IF EXISTS idx_magic_links_expires_at;
DROP INDEX IF EXISTS idx_magic_links_token;
DROP INDEX IF EXISTS idx_magic_links_user_id;

-- Drop tables
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS magic_links;

