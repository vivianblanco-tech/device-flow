-- Create shipment_type enum
CREATE TYPE shipment_type AS ENUM (
    'single_full_journey',
    'bulk_to_warehouse',
    'warehouse_to_engineer'
);

-- Add shipment_type column to shipments table
ALTER TABLE shipments ADD COLUMN shipment_type shipment_type;

-- Set default for existing shipments (backward compatibility)
UPDATE shipments SET shipment_type = 'single_full_journey' WHERE shipment_type IS NULL;

-- Make column NOT NULL after setting defaults
ALTER TABLE shipments ALTER COLUMN shipment_type SET NOT NULL;

-- Set default for new rows
ALTER TABLE shipments ALTER COLUMN shipment_type SET DEFAULT 'single_full_journey';

-- Add index for filtering by type
CREATE INDEX idx_shipments_type ON shipments(shipment_type);

-- Add comment
COMMENT ON COLUMN shipments.shipment_type IS 'Type of shipment: single_full_journey (1 laptop full flow), bulk_to_warehouse (2+ laptops to warehouse only), or warehouse_to_engineer (1 laptop from warehouse to engineer)';

