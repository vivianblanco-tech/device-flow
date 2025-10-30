-- Create shipment_status enum type
CREATE TYPE shipment_status AS ENUM (
    'pending_pickup',
    'picked_up_from_client',
    'in_transit_to_warehouse',
    'at_warehouse',
    'released_from_warehouse',
    'in_transit_to_engineer',
    'delivered'
);

-- Create shipments table
CREATE TABLE IF NOT EXISTS shipments (
    id BIGSERIAL PRIMARY KEY,
    client_company_id BIGINT NOT NULL REFERENCES client_companies(id) ON DELETE CASCADE,
    software_engineer_id BIGINT REFERENCES software_engineers(id) ON DELETE SET NULL,
    status shipment_status NOT NULL DEFAULT 'pending_pickup',
    courier_name VARCHAR(255),
    tracking_number VARCHAR(255),
    pickup_scheduled_date TIMESTAMP,
    picked_up_at TIMESTAMP,
    arrived_warehouse_at TIMESTAMP,
    released_warehouse_at TIMESTAMP,
    delivered_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_shipments_client_company_id ON shipments(client_company_id);
CREATE INDEX idx_shipments_software_engineer_id ON shipments(software_engineer_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
CREATE INDEX idx_shipments_created_at ON shipments(created_at);
CREATE INDEX idx_shipments_delivered_at ON shipments(delivered_at);

-- Comment on table and columns
COMMENT ON TABLE shipments IS 'Shipments tracking laptops through the delivery pipeline';
COMMENT ON COLUMN shipments.client_company_id IS 'Client company sending the laptops';
COMMENT ON COLUMN shipments.software_engineer_id IS 'Engineer receiving the laptop (optional, can be assigned later)';
COMMENT ON COLUMN shipments.status IS 'Current status in the delivery pipeline';
COMMENT ON COLUMN shipments.courier_name IS 'Name of courier service';
COMMENT ON COLUMN shipments.tracking_number IS 'Courier tracking number';
COMMENT ON COLUMN shipments.pickup_scheduled_date IS 'Scheduled pickup date from client';
COMMENT ON COLUMN shipments.picked_up_at IS 'When shipment was picked up from client';
COMMENT ON COLUMN shipments.arrived_warehouse_at IS 'When shipment arrived at warehouse';
COMMENT ON COLUMN shipments.released_warehouse_at IS 'When shipment was released from warehouse';
COMMENT ON COLUMN shipments.delivered_at IS 'When shipment was delivered to engineer';

