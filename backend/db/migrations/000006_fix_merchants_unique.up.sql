-- Drop old unique constraint and add new one
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS unique_merchant_name;
ALTER TABLE merchants ADD CONSTRAINT unique_merchant_name_family UNIQUE (name, family_id);
