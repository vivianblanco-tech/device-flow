-- Add employee_number column to software_engineers table
ALTER TABLE software_engineers
ADD COLUMN employee_number VARCHAR(50);

-- Add comment on column
COMMENT ON COLUMN software_engineers.employee_number IS 'Employee number (optional)';

