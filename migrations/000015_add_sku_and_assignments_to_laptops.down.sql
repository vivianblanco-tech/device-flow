-- Remove indexes
DROP INDEX IF EXISTS idx_laptops_software_engineer_id;
DROP INDEX IF EXISTS idx_laptops_client_company_id;
DROP INDEX IF EXISTS idx_laptops_sku;

-- Remove columns from laptops table
ALTER TABLE laptops DROP COLUMN IF EXISTS software_engineer_id;
ALTER TABLE laptops DROP COLUMN IF EXISTS client_company_id;
ALTER TABLE laptops DROP COLUMN IF EXISTS sku;

