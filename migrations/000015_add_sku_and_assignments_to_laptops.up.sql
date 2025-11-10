-- Add SKU, client company, and software engineer assignment fields to laptops table
ALTER TABLE laptops ADD COLUMN sku VARCHAR(100);
ALTER TABLE laptops ADD COLUMN client_company_id BIGINT REFERENCES client_companies(id) ON DELETE SET NULL;
ALTER TABLE laptops ADD COLUMN software_engineer_id BIGINT REFERENCES software_engineers(id) ON DELETE SET NULL;

-- Create indexes for better query performance
CREATE INDEX idx_laptops_sku ON laptops(sku);
CREATE INDEX idx_laptops_client_company_id ON laptops(client_company_id);
CREATE INDEX idx_laptops_software_engineer_id ON laptops(software_engineer_id);

-- Comment on new columns
COMMENT ON COLUMN laptops.sku IS 'Stock Keeping Unit - unique product identifier';
COMMENT ON COLUMN laptops.client_company_id IS 'Client company that owns/provided the laptop';
COMMENT ON COLUMN laptops.software_engineer_id IS 'Software engineer to whom the laptop is assigned';

