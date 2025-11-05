-- Revert shipment status from 'pending_pickup_from_client' back to 'pending_pickup'

-- Update all records back to the old value
UPDATE shipments SET status = 'pending_pickup' WHERE status = 'pending_pickup_from_client';

-- Note: Cannot remove enum value without recreating the type
-- Both values will exist in the enum type

