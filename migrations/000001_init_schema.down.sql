-- Rollback initial schema setup

-- Drop the version tracking table
DROP TABLE IF EXISTS schema_info;

-- Drop enum types
DROP TYPE IF EXISTS user_role;

-- Note: We don't drop the uuid-ossp extension as other databases might use it

