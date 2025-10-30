-- Create laptop_status enum type
CREATE TYPE laptop_status AS ENUM (
    'available',
    'in_transit_to_warehouse',
    'at_warehouse',
    'in_transit_to_engineer',
    'delivered',
    'retired'
);

-- Create laptops table
CREATE TABLE IF NOT EXISTS laptops (
    id BIGSERIAL PRIMARY KEY,
    serial_number VARCHAR(255) NOT NULL UNIQUE,
    brand VARCHAR(100),
    model VARCHAR(100),
    specs TEXT,
    status laptop_status NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_laptops_serial_number ON laptops(serial_number);
CREATE INDEX idx_laptops_status ON laptops(status);
CREATE INDEX idx_laptops_brand ON laptops(brand);

-- Create unique index on serial number (case-insensitive)
CREATE UNIQUE INDEX idx_laptops_serial_number_unique ON laptops(LOWER(serial_number));

-- Comment on table and columns
COMMENT ON TABLE laptops IS 'Laptop inventory tracking';
COMMENT ON COLUMN laptops.serial_number IS 'Unique serial number of the laptop';
COMMENT ON COLUMN laptops.brand IS 'Laptop brand (e.g., Dell, HP, Lenovo)';
COMMENT ON COLUMN laptops.model IS 'Laptop model';
COMMENT ON COLUMN laptops.specs IS 'Technical specifications';
COMMENT ON COLUMN laptops.status IS 'Current status in the delivery pipeline';

