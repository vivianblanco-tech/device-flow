-- Add client_company_id to users table for client role filtering
-- Note: This column may already exist from migration 000003
-- This migration ensures it exists with proper constraints

DO $$
BEGIN
    -- Add column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'client_company_id'
    ) THEN
        ALTER TABLE users ADD COLUMN client_company_id BIGINT REFERENCES client_companies(id) ON DELETE SET NULL;
    END IF;
    
    -- Create index if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE tablename = 'users' AND indexname = 'idx_users_client_company_id'
    ) THEN
        CREATE INDEX idx_users_client_company_id ON users(client_company_id);
    END IF;
END$$;

-- Add comment
COMMENT ON COLUMN users.client_company_id IS 'Client company ID (only for client role users) - used to filter shipments by company';

