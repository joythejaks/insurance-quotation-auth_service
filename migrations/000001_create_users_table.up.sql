-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "full_name" VARCHAR(255) NOT NULL,
    "role" VARCHAR(50) DEFAULT 'user',
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();