CREATE TABLE IF NOT EXISTS user(
    id serial PRIMARY KEY,
    email VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255)
);