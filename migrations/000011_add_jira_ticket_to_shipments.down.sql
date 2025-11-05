-- Remove index
DROP INDEX IF EXISTS idx_shipments_jira_ticket_number;

-- Remove jira_ticket_number column from shipments table
ALTER TABLE shipments 
DROP COLUMN IF EXISTS jira_ticket_number;

