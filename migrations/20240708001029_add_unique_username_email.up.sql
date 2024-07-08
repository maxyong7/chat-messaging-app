ALTER TABLE users
ADD CONSTRAINT unique_email_username UNIQUE (email, username);