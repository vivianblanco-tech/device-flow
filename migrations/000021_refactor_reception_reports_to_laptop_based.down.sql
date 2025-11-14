-- Rollback: Revert reception_reports from laptop-based back to shipment-based

-- Step 1: Drop the new reception_reports table
DROP TABLE IF EXISTS reception_reports;

-- Step 2: Recreate the old structure
CREATE TABLE IF NOT EXISTS reception_reports (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    warehouse_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    photo_urls TEXT[] DEFAULT ARRAY[]::TEXT[],
    expected_serial_number VARCHAR(255),
    actual_serial_number VARCHAR(255),
    serial_number_corrected BOOLEAN DEFAULT FALSE NOT NULL,
    correction_note TEXT,
    correction_approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);

-- Step 3: Recreate indexes
CREATE INDEX idx_reception_reports_shipment_id ON reception_reports(shipment_id);
CREATE INDEX idx_reception_reports_warehouse_user_id ON reception_reports(warehouse_user_id);
CREATE INDEX idx_reception_reports_received_at ON reception_reports(received_at);
CREATE INDEX idx_reception_reports_serial_corrected ON reception_reports(serial_number_corrected) WHERE serial_number_corrected = TRUE;

-- Step 4: Recreate unique constraint
CREATE UNIQUE INDEX idx_reception_reports_shipment_unique ON reception_reports(shipment_id);

-- Step 5: Add comments
COMMENT ON TABLE reception_reports IS 'Reports submitted by warehouse staff when receiving laptops';
COMMENT ON COLUMN reception_reports.photo_urls IS 'Array of photo URLs documenting received items';
COMMENT ON COLUMN reception_reports.expected_serial_number IS 'Serial number from pickup form (for single shipments)';
COMMENT ON COLUMN reception_reports.actual_serial_number IS 'Actual serial number received at warehouse';
COMMENT ON COLUMN reception_reports.serial_number_corrected IS 'Whether serial number was corrected from expected';
COMMENT ON COLUMN reception_reports.correction_note IS 'Note explaining why serial number was corrected';
COMMENT ON COLUMN reception_reports.correction_approved_by IS 'User ID (Logistics) who approved the correction';

-- Step 6: Drop the reception_report_status enum type
DROP TYPE IF EXISTS reception_report_status;

