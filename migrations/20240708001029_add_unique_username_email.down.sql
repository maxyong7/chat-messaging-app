ALTER TABLE users
REMOVE CONSTRAINT unique_email_username UNIQUE (email, username);