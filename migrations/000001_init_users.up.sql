CREATE TYPE user_types AS ENUM('Staff', 'Customer');

CREATE TABLE IF NOT EXISTS
users (
    id VARCHAR(16) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    user_type user_types NOT NULL,
    hashed_password BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS users_phone_number
	ON users(phone_number);
CREATE INDEX IF NOT EXISTS users_user_type
	ON users(user_type);
CREATE INDEX IF NOT EXISTS users_name
	ON users(lower(name));