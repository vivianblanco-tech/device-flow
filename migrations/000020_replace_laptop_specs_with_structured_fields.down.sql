-- Rollback: Restore specs column and make model nullable again

-- Add back the specs column
ALTER TABLE laptops ADD COLUMN specs TEXT;

-- Combine model, ram_gb, ssd_gb back into specs for existing rows
UPDATE laptops 
SET specs = CONCAT_WS(', ', 
    CASE WHEN model != '' THEN model ELSE NULL END,
    CASE WHEN ram_gb IS NOT NULL AND ram_gb != 'Unknown' THEN ram_gb || ' RAM' ELSE NULL END,
    CASE WHEN ssd_gb IS NOT NULL AND ssd_gb != 'Unknown' THEN ssd_gb || ' SSD' ELSE NULL END
);

-- Drop the new columns
ALTER TABLE laptops DROP COLUMN IF EXISTS ram_gb;
ALTER TABLE laptops DROP COLUMN IF EXISTS ssd_gb;

-- Make model nullable again
ALTER TABLE laptops ALTER COLUMN model DROP NOT NULL;

-- Restore comment
COMMENT ON COLUMN laptops.specs IS 'Technical specifications';

