-- Add CPU column to laptops table
ALTER TABLE laptops ADD COLUMN cpu TEXT NOT NULL DEFAULT '';

-- Update the constraint or add a check to ensure CPU is not empty for new records
-- For existing records, we set a default value of 'Unknown' temporarily
UPDATE laptops SET cpu = 'Unknown' WHERE cpu = '';

