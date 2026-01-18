-- Add unique constraint to merchants name
ALTER TABLE merchants ADD CONSTRAINT unique_merchant_name UNIQUE (name);
