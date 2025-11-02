-- Create Test User SQL Script
-- This creates a test user with email: admin@bairesdev.com and password: Test123!

-- Insert or update test user
-- Password hash is for "Test123!" (bcrypt cost 12)
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'admin@bairesdev.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS0MYq5IW',
    'logistics',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE
SET 
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW()
RETURNING id, email, role, created_at;

-- Display success message
\echo ''
\echo '==================================='
\echo 'Test User Created Successfully!'
\echo '==================================='
\echo 'Email:    admin@bairesdev.com'
\echo 'Password: Test123!'
\echo 'Role:     logistics'
\echo '==================================='
\echo ''
\echo 'Login at: http://localhost:8080/login'
\echo ''

