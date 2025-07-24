-- Add subadmin user to the database
-- Password: subadmin123 (hashed with bcrypt)

INSERT INTO users (username, password_hash, designation, is_admin, role, created_at) 
VALUES (
    'subadmin',
    '123456', -- password: subadmin123
    'Sub Administrator',
    false,
    'subadmin',
    NOW()
);

-- Verify the user was created
SELECT id, username, designation, is_admin, role, created_at FROM users WHERE username = 'subadmin'; 