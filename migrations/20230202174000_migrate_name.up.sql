CREATE TABLE IF NOT EXISTS links(
    id serial PRIMARY KEY,
    hash VARCHAR(255),
    to_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);