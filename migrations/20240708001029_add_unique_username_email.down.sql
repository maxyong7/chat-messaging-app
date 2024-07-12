ALTER TABLE IF EXISTS  users
REMOVE CONSTRAINT unique_email_username UNIQUE (email, username);