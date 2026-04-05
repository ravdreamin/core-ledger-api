-- Delete the dummy users and their records (cascaded by DB relations)
DELETE FROM users WHERE username IN ('admin', 'analyst', 'viewer');

-- Revert the records table alteration
ALTER TABLE records DROP COLUMN IF EXISTS category_id;
ALTER TABLE records ADD COLUMN IF NOT EXISTS category VARCHAR(100) NOT NULL DEFAULT 'Uncategorized';

-- Drop the categories table
DROP TABLE IF EXISTS categories;
