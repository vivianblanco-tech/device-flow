-- Add second_courier_name column to shipments table
ALTER TABLE shipments ADD COLUMN second_courier_name VARCHAR(255);

-- Create index for better query performance
CREATE INDEX idx_shipments_second_courier_name ON shipments(second_courier_name);

-- Add comment
COMMENT ON COLUMN shipments.second_courier_name IS 'Optional second courier name for shipments (e.g., return courier or secondary courier)';

