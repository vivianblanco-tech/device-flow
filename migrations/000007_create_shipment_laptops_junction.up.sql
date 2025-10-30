-- Create shipment_laptops junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS shipment_laptops (
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    laptop_id BIGINT NOT NULL REFERENCES laptops(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (shipment_id, laptop_id)
);

-- Create indexes for better query performance
CREATE INDEX idx_shipment_laptops_shipment_id ON shipment_laptops(shipment_id);
CREATE INDEX idx_shipment_laptops_laptop_id ON shipment_laptops(laptop_id);

-- Comment on table and columns
COMMENT ON TABLE shipment_laptops IS 'Junction table linking shipments to laptops (many-to-many)';
COMMENT ON COLUMN shipment_laptops.shipment_id IS 'Reference to shipment';
COMMENT ON COLUMN shipment_laptops.laptop_id IS 'Reference to laptop';

