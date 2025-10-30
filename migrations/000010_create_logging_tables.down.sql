-- Drop indexes for audit_logs
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_entity_id;
DROP INDEX IF EXISTS idx_audit_logs_entity_type;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;

-- Drop indexes for notification_logs
DROP INDEX IF EXISTS idx_notification_logs_status;
DROP INDEX IF EXISTS idx_notification_logs_sent_at;
DROP INDEX IF EXISTS idx_notification_logs_recipient;
DROP INDEX IF EXISTS idx_notification_logs_type;
DROP INDEX IF EXISTS idx_notification_logs_shipment_id;

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS notification_logs;

