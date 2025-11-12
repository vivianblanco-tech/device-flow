-- Remove index
DROP INDEX IF EXISTS idx_reception_reports_serial_corrected;

-- Remove columns
ALTER TABLE reception_reports DROP COLUMN IF EXISTS correction_approved_by;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS correction_note;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS serial_number_corrected;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS actual_serial_number;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS expected_serial_number;

