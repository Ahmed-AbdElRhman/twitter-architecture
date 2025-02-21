-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
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
    group_id INT  NOT NULL REFERENCES groups(id),
    endpoint_id INT  NOT NULL REFERENCES endpoints(id),
    UNIQUE (group_id, endpoint_id)
);

-- Create user_groups table
CREATE TABLE IF NOT EXISTS user_groups (
    id SERIAL PRIMARY KEY,
    user_id INT  NOT NULL REFERENCES users(id),
    group_id INT  NOT NULL REFERENCES groups(id)
    , UNIQUE (user_id, group_id)
);

-- Insert sample data into groups
INSERT INTO groups (name) VALUES ('admin'), ('user'),('manager'),('viewer'),('editor')
ON CONFLICT (name) DO NOTHING;

-- Insert sample data into endpoints
INSERT INTO endpoints (endpoint) VALUES ('view_tweets'), ('Create_tweets'), ('delete_tweets'),('edit_tweets')
ON CONFLICT (endpoint) DO NOTHING;


-- Insert sample data into users
INSERT INTO users (username, email, password) VALUES
('admin_user', 'admin@twt.com', '$2a$10$rE4z/c.EJ8Qv0gcRYoobjuZFn0RtfmtBFv5trdNm3Wc9Evvn1VBP6'), --password is P@ssw0rd
('normal_user', 'user@twt.com', '$2a$10$rE4z/c.EJ8Qv0gcRYoobjuZFn0RtfmtBFv5trdNm3Wc9Evvn1VBP6'), --password is P@ssw0rd
('manager_user', 'mngr@twt.com','$2a$10$rE4z/c.EJ8Qv0gcRYoobjuZFn0RtfmtBFv5trdNm3Wc9Evvn1VBP6') --password is P@ssw0rd
ON CONFLICT (email) DO NOTHING;

-- Insert sample data into user_groups
INSERT INTO user_groups (user_id, group_id)
SELECT (SELECT id FROM users WHERE email = 'admin@twt.com'), (SELECT id FROM groups WHERE name = 'admin')
UNION ALL
SELECT (SELECT id FROM users WHERE email = 'user@twt.com'), (SELECT id FROM groups WHERE name = 'viewer')
UNION ALL
SELECT (SELECT id FROM users WHERE email = 'user@twt.com'), (SELECT id FROM groups WHERE name = 'user')
UNION ALL
SELECT (SELECT id FROM users WHERE email = 'mngr@twt.com'), (SELECT id FROM groups WHERE name = 'manager')
UNION ALL
SELECT (SELECT id FROM users WHERE email = 'mngr@twt.com'), (SELECT id FROM groups WHERE name = 'editor')
WHERE NOT EXISTS (SELECT 1 FROM user_groups)
ON CONFLICT DO NOTHING;

-- Insert sample data into tweets only if the table is empty
INSERT INTO tweets (owner, tweet)
SELECT (SELECT id FROM users WHERE email = 'admin@twt.com'), 'This is a tweet by admin_user'
UNION ALL
SELECT(SELECT id FROM users WHERE email = 'user@twt.com'), 'This is a tweet by normal_user'
WHERE NOT EXISTS (SELECT 1 FROM tweets)
ON CONFLICT DO NOTHING;

-- Insert sample data into group_permissions only if the table is empty
INSERT INTO group_permissions (group_id, endpoint_id)
SELECT (SELECT id FROM groups WHERE name = 'user'), (SELECT id FROM endpoints WHERE endpoint = 'Create_tweets')
UNION ALL
SELECT (SELECT id FROM groups WHERE name = 'viewer'), (SELECT id FROM endpoints WHERE endpoint = 'view_tweets')
UNION ALL
SELECT (SELECT id FROM groups WHERE name = 'manager'), (SELECT id FROM endpoints WHERE endpoint = 'delete_tweets')
UNION ALL
SELECT (SELECT id FROM groups WHERE name = 'editor'), (SELECT id FROM endpoints WHERE endpoint = 'edit_tweets')
WHERE NOT EXISTS (SELECT 1 FROM group_permissions)
ON CONFLICT DO NOTHING;
----------- Create a view to get user data with endpoints -----------

-- Drop the view if it exists
DROP VIEW IF EXISTS user_data_with_endpoints;

-- Create the view
CREATE VIEW user_data_with_endpoints AS
SELECT 
    u.id as userid,
	g.name as groupname,
	ep.endpoint as authendpoint
FROM users u
inner JOIN user_groups ug ON u.id = ug.user_id
inner JOIN groups g ON ug.group_id = g.id
inner JOIN group_permissions gp ON g.id = gp.group_id
inner JOIN endpoints ep ON gp.endpoint_id = ep.id;