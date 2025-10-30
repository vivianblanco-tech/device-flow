-- Drop indexes
DROP INDEX IF EXISTS idx_laptops_serial_number_unique;
DROP INDEX IF EXISTS idx_laptops_brand;
DROP INDEX IF EXISTS idx_laptops_status;
DROP INDEX IF EXISTS idx_laptops_serial_number;

-- Drop laptops table
DROP TABLE IF EXISTS laptops;

-- Drop laptop_status enum type
DROP TYPE IF EXISTS laptop_status;

