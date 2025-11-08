-- Note: PostgreSQL does not support removing enum values directly.
-- To rollback, you would need to:
-- 1. Create a new enum without the value
-- 2. Convert the column to the new enum
-- 3. Drop the old enum
-- This is complex and typically not done in production.
-- For safety, this down migration is a no-op with a comment explaining the situation.

-- If you must remove the status, first ensure no rows use it:
-- UPDATE shipments SET status = 'pending_pickup_from_client' WHERE status = 'pickup_from_client_scheduled';

-- Then manually recreate the enum (requires dropping dependent objects):
-- This is intentionally left empty as it's a destructive operation.

