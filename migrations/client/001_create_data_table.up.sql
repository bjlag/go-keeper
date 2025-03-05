CREATE TABLE IF NOT EXISTS items (
    guid text PRIMARY KEY NOT NULL,
    categoryId INT NOT NULL,
    title text NOT NULL,
    value text,
    notes text NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);
