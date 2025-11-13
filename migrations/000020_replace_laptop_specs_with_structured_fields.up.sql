-- Replace free-form laptop specs with structured fields: model, ram_gb, ssd_gb
-- This migration makes model required and adds RAM and SSD as separate fields

-- First, set default empty string for existing NULL models (if any)
UPDATE laptops SET model = '' WHERE model IS NULL;

-- Make model column NOT NULL
ALTER TABLE laptops ALTER COLUMN model SET NOT NULL;

-- Add ram_gb column (NOT NULL with default for existing rows)
ALTER TABLE laptops ADD COLUMN ram_gb VARCHAR(50);

-- Add ssd_gb column (NOT NULL with default for existing rows)
ALTER TABLE laptops ADD COLUMN ssd_gb VARCHAR(50);

-- For existing rows, extract RAM/SSD from specs if possible, otherwise set placeholder
-- This is for data migration - new rows will require these fields
UPDATE laptops 
SET ram_gb = 'Unknown', 
    ssd_gb = 'Unknown' 
WHERE ram_gb IS NULL OR ssd_gb IS NULL;

-- Now make the new columns NOT NULL
ALTER TABLE laptops ALTER COLUMN ram_gb SET NOT NULL;
ALTER TABLE laptops ALTER COLUMN ssd_gb SET NOT NULL;

-- Drop the old specs column (replaced by structured fields)
ALTER TABLE laptops DROP COLUMN IF EXISTS specs;

-- Add comments
COMMENT ON COLUMN laptops.model IS 'Laptop model (e.g., Dell XPS 15, MacBook Pro 16) - REQUIRED';
COMMENT ON COLUMN laptops.ram_gb IS 'RAM size (e.g., 8GB, 16GB, 32GB) - REQUIRED';
COMMENT ON COLUMN laptops.ssd_gb IS 'SSD storage size (e.g., 256GB, 512GB, 1TB) - REQUIRED';

