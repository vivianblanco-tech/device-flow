-- Add new shipment status: 'pending_pickup_from_client'
-- Note: We keep 'pending_pickup' for backwards compatibility
-- New code should use 'pending_pickup_from_client'

ALTER TYPE shipment_status ADD VALUE IF NOT EXISTS 'pending_pickup_from_client';

