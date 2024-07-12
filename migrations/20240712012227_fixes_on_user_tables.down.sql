
ALTER TABLE user_credentials
DROP COLUMN user_uuid;

ALTER TABLE user_info DROP CONSTRAINT user_info_pkey;
ALTER TABLE user_info RENAME COLUMN user_uuid TO unique_id;
ALTER TABLE user_info DROP COLUMN id;
ALTER TABLE user_info ADD PRIMARY KEY (unique_id);