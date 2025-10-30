-- Remove foreign key from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_client_company;
DROP INDEX IF EXISTS idx_users_client_company_id;
ALTER TABLE users DROP COLUMN IF EXISTS client_company_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_client_companies_name_unique;
DROP INDEX IF EXISTS idx_client_companies_name;

-- Drop client_companies table
DROP TABLE IF EXISTS client_companies;

