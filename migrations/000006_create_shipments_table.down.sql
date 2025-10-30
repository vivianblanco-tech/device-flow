-- Drop indexes
DROP INDEX IF EXISTS idx_shipments_delivered_at;
DROP INDEX IF EXISTS idx_shipments_created_at;
DROP INDEX IF EXISTS idx_shipments_tracking_number;
DROP INDEX IF EXISTS idx_shipments_status;
DROP INDEX IF EXISTS idx_shipments_software_engineer_id;
DROP INDEX IF EXISTS idx_shipments_client_company_id;

-- Drop shipments table
DROP TABLE IF EXISTS shipments;

-- Drop shipment_status enum type
DROP TYPE IF EXISTS shipment_status;

