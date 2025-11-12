-- Add laptop_count column to shipments table
ALTER TABLE shipments ADD COLUMN laptop_count INTEGER;

-- Set default count for existing shipments based on shipment type
UPDATE shipments SET laptop_count = 1 WHERE shipment_type = 'single_full_journey' AND laptop_count IS NULL;
UPDATE shipments SET laptop_count = 1 WHERE shipment_type = 'warehouse_to_engineer' AND laptop_count IS NULL;

-- For bulk shipments, count existing laptops in junction table or set to a default
UPDATE shipments s
SET laptop_count = GREATEST((
    SELECT COUNT(*) FROM shipment_laptops sl WHERE sl.shipment_id = s.id
), 2)
WHERE s.shipment_type = 'bulk_to_warehouse' AND s.laptop_count IS NULL;

-- Make column NOT NULL after setting defaults
ALTER TABLE shipments ALTER COLUMN laptop_count SET NOT NULL;

-- Set default for new rows
ALTER TABLE shipments ALTER COLUMN laptop_count SET DEFAULT 1;

-- Add check constraint to ensure count is positive
ALTER TABLE shipments ADD CONSTRAINT chk_laptop_count_positive CHECK (laptop_count > 0);

-- Add comment
COMMENT ON COLUMN shipments.laptop_count IS 'Number of laptops in this shipment (must be 1 for single shipments, 2+ for bulk shipments)';

