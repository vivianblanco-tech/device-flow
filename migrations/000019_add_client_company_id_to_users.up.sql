-- Add client_company_id to users table for client role filtering
ALTER TABLE users ADD COLUMN client_company_id BIGINT REFERENCES client_companies(id) ON DELETE SET NULL;

-- Create index for better query performance
CREATE INDEX idx_users_client_company_id ON users(client_company_id);

-- Add comment
COMMENT ON COLUMN users.client_company_id IS 'Client company ID (only for client role users) - used to filter shipments by company';

