-- Rename shipment status from 'pending_pickup' to 'pending_pickup_from_client'

-- Step 1: Add new enum value
ALTER TYPE shipment_status ADD VALUE 'pending_pickup_from_client';

-- Step 2: Update all existing records to use the new value
UPDATE shipments SET status = 'pending_pickup_from_client' WHERE status = 'pending_pickup';

-- Step 3: Note - We cannot remove old enum values in PostgreSQL without recreating the type
-- This would require more complex migration with table recreation
-- For now, both values exist but only 'pending_pickup_from_client' should be used
-- The old 'pending_pickup' value is deprecated but remains for backwards compatibility

