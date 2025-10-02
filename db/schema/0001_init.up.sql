-- ========================================
-- 0001_init.up.sql
-- Nota CMS initial schema
-- ========================================

-- Enable pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================================
-- Function: update updated_at on row update
-- ========================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- Users table
-- ========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'editor', 'viewer')),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ========================================
-- Schemas table (user-defined models)
-- ========================================
CREATE TABLE schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    definition JSONB NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER trg_schemas_updated_at
BEFORE UPDATE ON schemas
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ========================================
-- Contents table
-- ========================================
CREATE TABLE contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id UUID REFERENCES schemas(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    published BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER trg_contents_updated_at
BEFORE UPDATE ON contents
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Index for quick access by content ID
CREATE INDEX idx_contents_id ON contents(id);

CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL,                         -- object key in S3
    url TEXT NOT NULL,                         -- presigned or public URL
    bucket TEXT NOT NULL,                      -- bucket name
    type TEXT NOT NULL CHECK (type IN ('image','video','file')),
    uploaded_by UUID REFERENCES users(id),     -- who uploaded it
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);
