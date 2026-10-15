CREATE TABLE api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  key_hash TEXT NOT NULL UNIQUE,   -- SHA-256 of the key, never the key itself
  prefix TEXT NOT NULL,            -- first 8 chars, for display in dashboard
  created_at TIMESTAMPTZ DEFAULT now(),
  revoked_at TIMESTAMPTZ
);
CREATE INDEX ON api_keys (tenant_id);