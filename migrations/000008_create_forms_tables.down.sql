-- Drop indexes for delivery_forms
DROP INDEX IF EXISTS idx_delivery_forms_shipment_unique;
DROP INDEX IF EXISTS idx_delivery_forms_delivered_at;
DROP INDEX IF EXISTS idx_delivery_forms_engineer_id;
DROP INDEX IF EXISTS idx_delivery_forms_shipment_id;

-- Drop indexes for reception_reports
DROP INDEX IF EXISTS idx_reception_reports_shipment_unique;
DROP INDEX IF EXISTS idx_reception_reports_received_at;
DROP INDEX IF EXISTS idx_reception_reports_warehouse_user_id;
DROP INDEX IF EXISTS idx_reception_reports_shipment_id;

-- Drop indexes for pickup_forms
DROP INDEX IF EXISTS idx_pickup_forms_shipment_unique;
DROP INDEX IF EXISTS idx_pickup_forms_submitted_at;
DROP INDEX IF EXISTS idx_pickup_forms_submitted_by_user_id;
DROP INDEX IF EXISTS idx_pickup_forms_shipment_id;

-- Drop tables
DROP TABLE IF EXISTS delivery_forms;
DROP TABLE IF EXISTS reception_reports;
DROP TABLE IF EXISTS pickup_forms;

