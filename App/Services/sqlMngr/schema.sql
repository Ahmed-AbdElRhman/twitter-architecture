-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    userid VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create tweets table
CREATE TABLE IF NOT EXISTS tweets (
    id SERIAL PRIMARY KEY,
    owner INT REFERENCES users(id),
    tweet TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Create endpoints table
CREATE TABLE IF NOT EXISTS endpoints (
    id SERIAL PRIMARY KEY,
    endpoint TEXT UNIQUE NOT NULL
);

-- Create group_permissions table
CREATE TABLE IF NOT EXISTS group_permissions (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id),
    endpoint_id INT REFERENCES endpoints(id)
);

-- Create user_groups table
CREATE TABLE IF NOT EXISTS user_groups (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    group_id INT REFERENCES groups(id)
);

-- Insert sample data into groups
INSERT INTO groups (name) VALUES ('admin'), ('user')
ON CONFLICT (name) DO NOTHING;

-- Insert sample data into endpoints
INSERT INTO endpoints (endpoint) VALUES ('view_tweets'), ('Create_tweets')
ON CONFLICT (endpoint) DO NOTHING;

-- Insert sample data into group_permissions
INSERT INTO group_permissions (group_id, endpoint_id) VALUES
((SELECT id FROM groups WHERE name = 'admin'), (SELECT id FROM endpoints WHERE endpoint = 'view_products')),
((SELECT id FROM groups WHERE name = 'user'), (SELECT id FROM endpoints WHERE endpoint = 'view_profile'))
ON CONFLICT DO NOTHING;

-- Insert sample data into users
INSERT INTO users (userid, email, password) VALUES
('admin_user', 'admin@twt.com', 'admin'),
('normal_user', 'user@twt.com', 'user')
ON CONFLICT (email) DO NOTHING;

-- Insert sample data into user_groups
INSERT INTO user_groups (user_id, group_id) VALUES
((SELECT id FROM users WHERE email = 'admin@twt.com'), (SELECT id FROM groups WHERE name = 'admin')),
((SELECT id FROM users WHERE email = 'user@twt.com'), (SELECT id FROM groups WHERE name = 'user'))
ON CONFLICT DO NOTHING;

-- Insert sample data into tweets
INSERT INTO tweets (owner, tweet) VALUES
((SELECT id FROM users WHERE email = 'admin@twt.com'), 'This is a tweet by admin_user'),
((SELECT id FROM users WHERE email = 'user@twt.com'), 'This is a tweet by normal_user')
ON CONFLICT DO NOTHING;