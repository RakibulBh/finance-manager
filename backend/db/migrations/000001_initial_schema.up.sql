-- Create Extension for UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Families Table
CREATE TABLE families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    currency TEXT DEFAULT 'USD' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    email TEXT UNIQUE NOT NULL,
    password_digest TEXT NOT NULL,
    role TEXT DEFAULT 'member' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Accounts Table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance DECIMAL(19,4) DEFAULT 0 NOT NULL,
    currency TEXT NOT NULL,
    subtype TEXT, -- e.g., 'checking', 'savings', 'credit_card'
    classification TEXT NOT NULL, -- 'asset' or 'liability'
    status TEXT DEFAULT 'active' NOT NULL,
    plaid_account_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Entries Table (The Ledger)
CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount DECIMAL(19,4) NOT NULL,
    date DATE NOT NULL,
    currency TEXT NOT NULL,
    name TEXT NOT NULL,
    entryable_type TEXT NOT NULL, -- 'Transaction', 'Valuation', 'Trade'
    entryable_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Transactions Table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID, -- To be linked to categories table later
    merchant_id UUID, -- To be linked to merchants table later
    kind TEXT DEFAULT 'standard' NOT NULL
);

-- Indices for Performance
CREATE INDEX idx_users_family_id ON users(family_id);
CREATE INDEX idx_accounts_family_id ON accounts(family_id);
CREATE INDEX idx_entries_account_id ON entries(account_id);
CREATE INDEX idx_entries_date ON entries(date);
CREATE INDEX idx_entries_entryable ON entries(entryable_type, entryable_id);
