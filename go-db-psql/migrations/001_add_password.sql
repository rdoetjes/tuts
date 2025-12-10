-- Bootstrap migration: Add password column to users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(128),
    lastname VARCHAR(128),
    dob DATE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255)
);

-- Check and add password column (inittially it didn't exist)
ALTER TABLE users
ADD COLUMN IF NOT EXISTS password VARCHAR(255) NOT NULL DEFAULT '';

-- Create an initial admin user with a default password (CHANGE THIS IN PRODUCTION)
-- Password: 'admin123' (hashed with bcrypt)
-- To generate your own: use the setup.go tool or an online bcrypt generator
INSERT INTO users (firstname, lastname, email, password, dob)
VALUES ('Admin', 'User', 'admin@localhost', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36gBS/Pq', '2000-01-01')
ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password;
