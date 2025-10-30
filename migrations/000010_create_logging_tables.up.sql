-- Create notification_logs table
CREATE TABLE IF NOT EXISTS notification_logs (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT REFERENCES shipments(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    details JSONB DEFAULT '{}'::jsonb
);

-- Create indexes for better query performance
CREATE INDEX idx_notification_logs_shipment_id ON notification_logs(shipment_id);
CREATE INDEX idx_notification_logs_type ON notification_logs(type);
CREATE INDEX idx_notification_logs_recipient ON notification_logs(recipient);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);

-- Create composite index for entity lookups
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

-- Comment on tables and columns
COMMENT ON TABLE notification_logs IS 'Log of all notifications sent by the system';
COMMENT ON COLUMN notification_logs.type IS 'Type of notification (e.g., email, sms)';
COMMENT ON COLUMN notification_logs.recipient IS 'Recipient email address or phone number';
COMMENT ON COLUMN notification_logs.status IS 'Status of notification (sent, failed, pending)';

COMMENT ON TABLE audit_logs IS 'Audit trail of important system actions';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (create, update, delete, etc.)';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity affected (shipment, laptop, user, etc.)';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the affected entity';
COMMENT ON COLUMN audit_logs.details IS 'Additional details about the action (JSON)';

