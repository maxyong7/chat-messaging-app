ALTER TABLE IF EXISTS users
RENAME TO user_credentials;

CREATE TABLE IF NOT EXISTS user_info (
    unique_id TEXT DEFAULT gen_random_uuid (),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    avatar TEXT,
    PRIMARY KEY (unique_id)
);

CREATE TABLE IF NOT EXISTS contacts (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_uuid TEXT,
    contact_user_uuid TEXT,
    conversation_uuid TEXT DEFAULT gen_random_uuid (),
    blocked BOOLEAN
);


CREATE TABLE IF NOT EXISTS participants (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_uuid TEXT,
    conversation_uuid TEXT,
    join_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    left_date TIMESTAMPTZ
); 

CREATE TABLE IF NOT EXISTS conversations (
    conversation_uuid TEXT DEFAULT gen_random_uuid (),
    last_message TEXT,
    last_sent_user_uuid TEXT,
    title TEXT,
    last_message_created_at TIMESTAMPTZ,
    conversation_type TEXT,
    PRIMARY KEY (conversation_uuid)
); 

CREATE TABLE IF NOT EXISTS messages (
    message_uuid TEXT DEFAULT gen_random_uuid (),
    conversation_uuid TEXT,
    user_uuid TEXT,
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN,
    PRIMARY KEY (message_uuid)
); 

CREATE TABLE IF NOT EXISTS seen_status (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    message_uuid TEXT,
    user_uuid TEXT,
    seen_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE IF NOT EXISTS reaction (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    message_uuid TEXT,
    user_uuid TEXT,
    reaction_type TEXT
); 