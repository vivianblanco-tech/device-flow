-- Add jira_ticket_number column to shipments table
ALTER TABLE shipments 
ADD COLUMN jira_ticket_number VARCHAR(50) NOT NULL DEFAULT '';

-- Remove the default after adding the column
ALTER TABLE shipments 
ALTER COLUMN jira_ticket_number DROP DEFAULT;

-- Add index for better query performance on JIRA ticket lookups
CREATE INDEX idx_shipments_jira_ticket_number ON shipments(jira_ticket_number);

-- Add comment on column
COMMENT ON COLUMN shipments.jira_ticket_number IS 'JIRA ticket number associated with this shipment (format: PROJECT-NUMBER)';

