
ALTER TABLE IF EXISTS user_credentials
DROP COLUMN user_uuid;

ALTER TABLE IF EXISTS user_info DROP CONSTRAINT user_info_pkey;
ALTER TABLE IF EXISTS user_info RENAME COLUMN user_uuid TO unique_id;
ALTER TABLE IF EXISTS user_info DROP COLUMN id;
ALTER TABLE IF EXISTS user_info ADD PRIMARY KEY (unique_id);