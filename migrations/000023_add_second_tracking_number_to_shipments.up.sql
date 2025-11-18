-- Add second_tracking_number column to shipments table
ALTER TABLE shipments ADD COLUMN second_tracking_number VARCHAR(255);

-- Add index for better query performance
CREATE INDEX idx_shipments_second_tracking_number ON shipments(second_tracking_number);

-- Add comment
COMMENT ON COLUMN shipments.second_tracking_number IS 'Optional second tracking number for shipments (e.g., return tracking or secondary courier)';

