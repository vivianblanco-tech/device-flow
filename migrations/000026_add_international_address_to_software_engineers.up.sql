-- Add international address fields to software_engineers table
ALTER TABLE software_engineers
ADD COLUMN IF NOT EXISTS address_street VARCHAR(255),
ADD COLUMN IF NOT EXISTS address_city VARCHAR(255),
ADD COLUMN IF NOT EXISTS address_country VARCHAR(255),
ADD COLUMN IF NOT EXISTS address_state VARCHAR(255),
ADD COLUMN IF NOT EXISTS address_postal_code VARCHAR(50);

-- Add comments
COMMENT ON COLUMN software_engineers.address_street IS 'Street address (international format)';
COMMENT ON COLUMN software_engineers.address_city IS 'City (international format)';
COMMENT ON COLUMN software_engineers.address_country IS 'Country (international format)';
COMMENT ON COLUMN software_engineers.address_state IS 'State/Province (optional, international format)';
COMMENT ON COLUMN software_engineers.address_postal_code IS 'Postal/ZIP code (optional, international format)';

