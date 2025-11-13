-- Remove index
DROP INDEX IF EXISTS idx_users_client_company_id;

-- Remove column from users table
ALTER TABLE users DROP COLUMN IF EXISTS client_company_id;

