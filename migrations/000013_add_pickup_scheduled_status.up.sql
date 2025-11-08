-- Add pickup_from_client_scheduled status to shipment_status enum
ALTER TYPE shipment_status ADD VALUE IF NOT EXISTS 'pickup_from_client_scheduled' BEFORE 'picked_up_from_client';

