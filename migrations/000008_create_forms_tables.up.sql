-- Create pickup_forms table
CREATE TABLE IF NOT EXISTS pickup_forms (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    submitted_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    form_data JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- Create reception_reports table
CREATE TABLE IF NOT EXISTS reception_reports (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    warehouse_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    photo_urls TEXT[] DEFAULT ARRAY[]::TEXT[]
);

-- Create delivery_forms table
CREATE TABLE IF NOT EXISTS delivery_forms (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    engineer_id BIGINT NOT NULL REFERENCES software_engineers(id) ON DELETE CASCADE,
    delivered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    photo_urls TEXT[] DEFAULT ARRAY[]::TEXT[]
);

-- Create indexes for better query performance
CREATE INDEX idx_pickup_forms_shipment_id ON pickup_forms(shipment_id);
CREATE INDEX idx_pickup_forms_submitted_by_user_id ON pickup_forms(submitted_by_user_id);
CREATE INDEX idx_pickup_forms_submitted_at ON pickup_forms(submitted_at);

CREATE INDEX idx_reception_reports_shipment_id ON reception_reports(shipment_id);
CREATE INDEX idx_reception_reports_warehouse_user_id ON reception_reports(warehouse_user_id);
CREATE INDEX idx_reception_reports_received_at ON reception_reports(received_at);

CREATE INDEX idx_delivery_forms_shipment_id ON delivery_forms(shipment_id);
CREATE INDEX idx_delivery_forms_engineer_id ON delivery_forms(engineer_id);
CREATE INDEX idx_delivery_forms_delivered_at ON delivery_forms(delivered_at);

-- Add unique constraints to ensure one form per shipment
CREATE UNIQUE INDEX idx_pickup_forms_shipment_unique ON pickup_forms(shipment_id);
CREATE UNIQUE INDEX idx_reception_reports_shipment_unique ON reception_reports(shipment_id);
CREATE UNIQUE INDEX idx_delivery_forms_shipment_unique ON delivery_forms(shipment_id);

-- Comment on tables and columns
COMMENT ON TABLE pickup_forms IS 'Forms submitted by clients to schedule laptop pickup';
COMMENT ON COLUMN pickup_forms.form_data IS 'JSON data containing form fields (flexible schema)';

COMMENT ON TABLE reception_reports IS 'Reports submitted by warehouse staff when receiving laptops';
COMMENT ON COLUMN reception_reports.photo_urls IS 'Array of photo URLs documenting received items';

COMMENT ON TABLE delivery_forms IS 'Forms submitted when laptops are delivered to engineers';
COMMENT ON COLUMN delivery_forms.photo_urls IS 'Array of photo URLs documenting delivery';

