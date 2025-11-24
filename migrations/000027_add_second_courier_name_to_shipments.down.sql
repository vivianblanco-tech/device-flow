-- Remove index
DROP INDEX IF EXISTS idx_shipments_second_courier_name;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS second_courier_name;

