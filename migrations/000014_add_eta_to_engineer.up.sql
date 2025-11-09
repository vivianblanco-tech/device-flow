-- Add eta_to_engineer column to shipments table
ALTER TABLE shipments ADD COLUMN eta_to_engineer TIMESTAMP;

-- Add index for better query performance
CREATE INDEX idx_shipments_eta_to_engineer ON shipments(eta_to_engineer);

-- Comment on column
COMMENT ON COLUMN shipments.eta_to_engineer IS 'Estimated time of arrival to engineer when status is in_transit_to_engineer';

