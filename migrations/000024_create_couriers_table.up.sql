-- Create couriers table
CREATE TABLE IF NOT EXISTS couriers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    contact_info TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_couriers_name ON couriers(name);

-- Comment on table and columns
COMMENT ON TABLE couriers IS 'Courier companies used for shipments';
COMMENT ON COLUMN couriers.name IS 'Courier company name (unique)';
COMMENT ON COLUMN couriers.contact_info IS 'Contact information for the courier company';

