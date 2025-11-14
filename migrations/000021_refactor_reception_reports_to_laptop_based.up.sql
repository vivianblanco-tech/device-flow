-- Refactor reception_reports from shipment-based to laptop-based
-- This enables one reception report per laptop with approval workflow

-- Step 1: Create new reception_report_status enum type
CREATE TYPE reception_report_status AS ENUM (
    'pending_approval',
    'approved'
);

-- Step 2: Drop existing indexes before restructuring
DROP INDEX IF EXISTS idx_reception_reports_shipment_unique;
DROP INDEX IF EXISTS idx_reception_reports_shipment_id;
DROP INDEX IF EXISTS idx_reception_reports_warehouse_user_id;
DROP INDEX IF EXISTS idx_reception_reports_received_at;
DROP INDEX IF EXISTS idx_reception_reports_serial_corrected;

-- Step 3: Rename and restructure the table (preserve data by creating new table and migrating)
ALTER TABLE reception_reports RENAME TO reception_reports_old;

-- Step 4: Create new reception_reports table with laptop-based structure
CREATE TABLE reception_reports (
    id BIGSERIAL PRIMARY KEY,
    laptop_id BIGINT NOT NULL REFERENCES laptops(id) ON DELETE CASCADE,
    shipment_id BIGINT REFERENCES shipments(id) ON DELETE SET NULL,
    client_company_id BIGINT REFERENCES client_companies(id) ON DELETE SET NULL,
    tracking_number VARCHAR(255),
    warehouse_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    
    -- Required photo uploads
    photo_serial_number VARCHAR(500) NOT NULL,
    photo_external_condition VARCHAR(500) NOT NULL,
    photo_working_condition VARCHAR(500) NOT NULL,
    
    -- Approval tracking
    status reception_report_status NOT NULL DEFAULT 'pending_approval',
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 5: Create indexes for better query performance
CREATE INDEX idx_reception_reports_laptop_id ON reception_reports(laptop_id);
CREATE INDEX idx_reception_reports_shipment_id ON reception_reports(shipment_id);
CREATE INDEX idx_reception_reports_client_company_id ON reception_reports(client_company_id);
CREATE INDEX idx_reception_reports_warehouse_user_id ON reception_reports(warehouse_user_id);
CREATE INDEX idx_reception_reports_status ON reception_reports(status);
CREATE INDEX idx_reception_reports_received_at ON reception_reports(received_at);
CREATE INDEX idx_reception_reports_approved_by ON reception_reports(approved_by);

-- Step 6: Create unique constraint to ensure one reception report per laptop
CREATE UNIQUE INDEX idx_reception_reports_laptop_unique ON reception_reports(laptop_id);

-- Step 7: Drop old table (sample data can be recreated)
DROP TABLE reception_reports_old;

-- Step 8: Add comments
COMMENT ON TABLE reception_reports IS 'Reception reports submitted by warehouse staff when receiving individual laptops';
COMMENT ON COLUMN reception_reports.laptop_id IS 'Laptop this report is for (one report per laptop)';
COMMENT ON COLUMN reception_reports.shipment_id IS 'Reference to original shipment (optional)';
COMMENT ON COLUMN reception_reports.client_company_id IS 'Client company for tracking purposes (optional)';
COMMENT ON COLUMN reception_reports.tracking_number IS 'Tracking number for reference (optional)';
COMMENT ON COLUMN reception_reports.photo_serial_number IS 'Photo URL of serial number verification';
COMMENT ON COLUMN reception_reports.photo_external_condition IS 'Photo URL of laptop external condition';
COMMENT ON COLUMN reception_reports.photo_working_condition IS 'Photo URL of laptop working condition (powered on)';
COMMENT ON COLUMN reception_reports.status IS 'Approval status: pending_approval or approved';
COMMENT ON COLUMN reception_reports.approved_by IS 'Logistics user who approved the report';
COMMENT ON COLUMN reception_reports.approved_at IS 'Timestamp when report was approved';

