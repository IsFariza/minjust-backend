CREATE TABLE IF NOT EXISTS admin_accounts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO admin_accounts (username, password_hash) 
VALUES ('admin', '$2a$10$1GNjZzXsGMxCgu6FEwvTsuaz6VKmZeT9Hc10h4QRI3c/D08hoHy0W')
ON CONFLICT (username) DO NOTHING;