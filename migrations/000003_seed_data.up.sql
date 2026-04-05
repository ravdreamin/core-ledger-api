CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Note: Modifying records table to support the category JOIN
ALTER TABLE records DROP COLUMN IF EXISTS category;
ALTER TABLE records ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- Insert 5 categories
INSERT INTO categories (name) VALUES 
    ('Rent'), 
    ('Food'), 
    ('Salary'), 
    ('Utilities'), 
    ('Entertainment')
ON CONFLICT DO NOTHING;

-- Insert 3 users
INSERT INTO users (username, password_hash, role) VALUES 
    ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
    ('analyst', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'analyst'),
    ('viewer', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer')
ON CONFLICT DO NOTHING;

-- We need to fetch the actual IDs assigned. Alternatively, use a CTE for inserts
WITH user_ids AS (
    SELECT id, username FROM users WHERE username IN ('admin', 'analyst', 'viewer')
), category_ids AS (
    SELECT id, name FROM categories WHERE name IN ('Rent', 'Food', 'Salary', 'Utilities', 'Entertainment')
)
INSERT INTO records (user_id, amount, type, date, description, category_id)
SELECT u.id, 500, 'expense', NOW(), 'Monthly rent', c.id
FROM user_ids u CROSS JOIN category_ids c WHERE c.name = 'Rent' AND u.username = 'admin'
UNION ALL
SELECT u.id, 100, 'expense', NOW(), 'Groceries', c.id
FROM user_ids u CROSS JOIN category_ids c WHERE c.name = 'Food' AND u.username = 'admin'
UNION ALL
