-- Initial schema setup
-- This migration creates the foundation for Align

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types that will be used across tables
CREATE TYPE user_role AS ENUM ('logistics', 'client', 'warehouse', 'project_manager');

-- Create a simple version tracking table to verify migrations work
CREATE TABLE IF NOT EXISTS schema_info (
    id SERIAL PRIMARY KEY,
    version VARCHAR(50) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Insert initial version
INSERT INTO schema_info (version, description) 
VALUES ('000001', 'Initial schema setup with enums and version tracking');

-- Add comment
COMMENT ON TABLE schema_info IS 'Tracks schema versions and migration history';

