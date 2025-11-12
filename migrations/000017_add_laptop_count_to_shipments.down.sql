-- Remove check constraint
ALTER TABLE shipments DROP CONSTRAINT IF EXISTS chk_laptop_count_positive;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS laptop_count;

