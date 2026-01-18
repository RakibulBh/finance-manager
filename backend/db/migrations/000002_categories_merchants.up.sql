-- Categories Table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    color TEXT DEFAULT '#6172F3' NOT NULL,
    classification TEXT DEFAULT 'expense' NOT NULL, -- 'income' or 'expense'
    lucide_icon TEXT DEFAULT 'shapes' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Merchants Table
CREATE TABLE merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID REFERENCES families(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT,
    logo_url TEXT,
    website_url TEXT,
    source TEXT, -- e.g., 'manual', 'plaid'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Add Foreign Keys to Transactions (Now that tables exist)
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;

ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE SET NULL;

-- Indices
CREATE INDEX idx_categories_family_id ON categories(family_id);
CREATE INDEX idx_merchants_family_id ON merchants(family_id);
