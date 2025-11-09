-- Remove index
DROP INDEX IF EXISTS idx_shipments_eta_to_engineer;

-- Remove eta_to_engineer column from shipments table
ALTER TABLE shipments DROP COLUMN IF EXISTS eta_to_engineer;

