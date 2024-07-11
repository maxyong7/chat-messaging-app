
ALTER TABLE user_credentials 
RENAME TO users;

DROP TABLE IF EXISTS user_info;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS seen_status;
DROP TABLE IF EXISTS reaction;