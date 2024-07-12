ALTER TABLE IF EXISTS  users
ADD CONSTRAINT unique_email_username UNIQUE (email, username);