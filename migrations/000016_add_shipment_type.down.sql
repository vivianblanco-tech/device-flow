-- Remove index
DROP INDEX IF EXISTS idx_shipments_type;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS shipment_type;

-- Drop enum type
DROP TYPE IF EXISTS shipment_type;

