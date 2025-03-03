CREATE TABLE IF NOT EXISTS data (
    guid           text PRIMARY KEY,
    encrypted_data text NOT NULL,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);
