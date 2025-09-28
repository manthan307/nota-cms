-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'editor', 'viewer')),
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP NULL
);

-- Schemas (user-defined models)
CREATE TABLE schemas (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    definition JSONB NOT NULL,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP NULL
);

-- Contents
CREATE TABLE contents (
    id SERIAL PRIMARY KEY,
    schema_id INT REFERENCES schemas(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP NULL
);
