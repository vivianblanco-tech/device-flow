-- Remove employee_number column from software_engineers table
ALTER TABLE software_engineers
DROP COLUMN IF EXISTS employee_number;

