INSERT INTO auth_users (name, password, email)
VALUES ('admin', 'admin', 'admin@localhost') ON DUPLICATE KEY
UPDATE email = email;