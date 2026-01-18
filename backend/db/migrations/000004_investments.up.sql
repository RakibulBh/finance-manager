-- Securities Table
CREATE TABLE securities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    latest_price DECIMAL(19,4) DEFAULT 0 NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Trades Table
CREATE TABLE trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Fixed ID for polymorphic link? No, Entries(entryable_id) links here.
    security_id UUID NOT NULL REFERENCES securities(id) ON DELETE CASCADE,
    qty DECIMAL(19,4) NOT NULL,
    price DECIMAL(19,4) NOT NULL,
    kind TEXT DEFAULT 'buy' NOT NULL -- 'buy' or 'sell'
);

-- Indices
CREATE INDEX idx_securities_ticker ON securities(ticker);
CREATE INDEX idx_trades_security_id ON trades(security_id);
