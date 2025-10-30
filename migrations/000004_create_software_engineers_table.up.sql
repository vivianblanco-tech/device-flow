-- Create software_engineers table
CREATE TABLE IF NOT EXISTS software_engineers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    address_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    address_confirmation_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_software_engineers_email ON software_engineers(email);
CREATE INDEX idx_software_engineers_name ON software_engineers(name);
CREATE INDEX idx_software_engineers_address_confirmed ON software_engineers(address_confirmed);

-- Add unique constraint on email
CREATE UNIQUE INDEX idx_software_engineers_email_unique ON software_engineers(LOWER(email));

-- Comment on table and columns
COMMENT ON TABLE software_engineers IS 'Software engineers who receive laptops';
COMMENT ON COLUMN software_engineers.name IS 'Full name of the software engineer';
COMMENT ON COLUMN software_engineers.email IS 'Email address (unique, case-insensitive)';
COMMENT ON COLUMN software_engineers.address IS 'Delivery address';
COMMENT ON COLUMN software_engineers.phone IS 'Phone number';
COMMENT ON COLUMN software_engineers.address_confirmed IS 'Whether the engineer has confirmed their address';
COMMENT ON COLUMN software_engineers.address_confirmation_at IS 'Timestamp when address was confirmed';

