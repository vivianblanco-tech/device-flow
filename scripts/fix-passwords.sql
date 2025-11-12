-- Fix all user passwords to Test123!
-- Password hash: $2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK

UPDATE users 
SET password_hash = '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK';

-- Verify the update
SELECT id, email, role, LEFT(password_hash, 20) as hash_start 
FROM users 
ORDER BY id 
LIMIT 5;

