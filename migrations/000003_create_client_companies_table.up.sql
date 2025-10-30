-- Create client_companies table
CREATE TABLE IF NOT EXISTS client_companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_info TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_client_companies_name ON client_companies(name);

-- Add unique constraint on company name to avoid duplicates
CREATE UNIQUE INDEX idx_client_companies_name_unique ON client_companies(LOWER(name));

-- Add foreign key to users table to link users to client companies
ALTER TABLE users ADD COLUMN client_company_id BIGINT;
ALTER TABLE users ADD CONSTRAINT fk_users_client_company 
    FOREIGN KEY (client_company_id) 
    REFERENCES client_companies(id) 
    ON DELETE SET NULL;

-- Create index for the foreign key
CREATE INDEX idx_users_client_company_id ON users(client_company_id);

-- Comment on table and columns
COMMENT ON TABLE client_companies IS 'Client companies that send laptops for pickup';
COMMENT ON COLUMN client_companies.name IS 'Company name (unique, case-insensitive)';
COMMENT ON COLUMN client_companies.contact_info IS 'Contact information (email, phone, etc.)';
COMMENT ON COLUMN users.client_company_id IS 'Reference to client company (for client role users)';

