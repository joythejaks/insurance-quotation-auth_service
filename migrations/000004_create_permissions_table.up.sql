CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code);
CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();