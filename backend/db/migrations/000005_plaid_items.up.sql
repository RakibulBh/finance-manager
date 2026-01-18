-- Plaid Items Table
CREATE TABLE plaid_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    item_id TEXT UNIQUE NOT NULL,
    institution_id TEXT,
    institution_name TEXT,
    sync_cursor TEXT,
    status TEXT DEFAULT 'active' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Index for family lookup
CREATE INDEX idx_plaid_items_family_id ON plaid_items(family_id);
