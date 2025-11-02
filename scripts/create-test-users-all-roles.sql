-- Create Test Users for All Roles
-- This creates test users for each role in the system
-- All users have the password: Test123!
-- Password hash: $2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK

-- Logistics user (full access)
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'logistics@bairesdev.com',
    '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK',
    'logistics',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE
SET 
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW();

-- Client user (limited access)
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'client@bairesdev.com',
    '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK',
    'client',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE
SET 
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW();

-- Warehouse user (medium access)
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'warehouse@bairesdev.com',
    '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK',
    'warehouse',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE
SET 
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW();

-- Project Manager user (read-only access)
INSERT INTO users (email, password_hash, role, created_at, updated_at)
VALUES (
    'pm@bairesdev.com',
    '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK',
    'project_manager',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE
SET 
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    updated_at = NOW();

-- Display results
\echo ''
\echo '======================================='
\echo 'Test Users Created Successfully!'
\echo '======================================='
\echo ''
\echo 'All users have password: Test123!'
\echo ''
\echo 'LOGISTICS USER (Full Access):'
\echo '  Email: logistics@bairesdev.com'
\echo '  Role:  logistics'
\echo ''
\echo 'CLIENT USER (Limited Access):'
\echo '  Email: client@bairesdev.com'
\echo '  Role:  client'
\echo ''
\echo 'WAREHOUSE USER (Medium Access):'
\echo '  Email: warehouse@bairesdev.com'
\echo '  Role:  warehouse'
\echo ''
\echo 'PROJECT MANAGER USER (Read-Only):'
\echo '  Email: pm@bairesdev.com'
\echo '  Role:  project_manager'
\echo ''
\echo '======================================='
\echo ''
\echo 'Login at: http://localhost:8080/login'
\echo ''

-- Show all test users
SELECT id, email, role, created_at 
FROM users 
WHERE email IN (
    'logistics@bairesdev.com',
    'client@bairesdev.com',
    'warehouse@bairesdev.com',
    'pm@bairesdev.com'
)
ORDER BY 
    CASE role
        WHEN 'logistics' THEN 1
        WHEN 'warehouse' THEN 2
        WHEN 'project_manager' THEN 3
        WHEN 'client' THEN 4
    END;

