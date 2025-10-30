-- Drop indexes
DROP INDEX IF EXISTS idx_shipment_laptops_laptop_id;
DROP INDEX IF EXISTS idx_shipment_laptops_shipment_id;

-- Drop shipment_laptops junction table
DROP TABLE IF EXISTS shipment_laptops;

