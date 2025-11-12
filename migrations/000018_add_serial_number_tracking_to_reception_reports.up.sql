-- Add serial number tracking columns to reception_reports
ALTER TABLE reception_reports ADD COLUMN expected_serial_number VARCHAR(255);
ALTER TABLE reception_reports ADD COLUMN actual_serial_number VARCHAR(255);
ALTER TABLE reception_reports ADD COLUMN serial_number_corrected BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE reception_reports ADD COLUMN correction_note TEXT;
ALTER TABLE reception_reports ADD COLUMN correction_approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL;

-- Add index for finding corrected serial numbers
CREATE INDEX idx_reception_reports_serial_corrected ON reception_reports(serial_number_corrected) WHERE serial_number_corrected = TRUE;

-- Add comments
COMMENT ON COLUMN reception_reports.expected_serial_number IS 'Serial number from pickup form (for single shipments)';
COMMENT ON COLUMN reception_reports.actual_serial_number IS 'Actual serial number received at warehouse';
COMMENT ON COLUMN reception_reports.serial_number_corrected IS 'Whether serial number was corrected from expected';
COMMENT ON COLUMN reception_reports.correction_note IS 'Note explaining why serial number was corrected';
COMMENT ON COLUMN reception_reports.correction_approved_by IS 'User ID (Logistics) who approved the correction';

