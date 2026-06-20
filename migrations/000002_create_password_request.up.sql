CREATE TABLE IF NOT EXISTS password_requests (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employee_accounts(id),
    system_name VARCHAR(50) NOT NULL,
    input_data JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    primary_password TEXT,
    admin_comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
