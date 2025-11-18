-- Remove index
DROP INDEX IF EXISTS idx_shipments_second_tracking_number;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS second_tracking_number;

