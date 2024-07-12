ALTER TABLE user_credentials
ADD user_uuid TEXT DEFAULT gen_random_uuid (); 

ALTER TABLE user_info DROP CONSTRAINT user_info_pkey;
ALTER TABLE user_info RENAME COLUMN unique_id TO user_uuid;
ALTER TABLE user_info ADD id BIGINT;
ALTER TABLE user_info ADD PRIMARY KEY (id);
ALTER TABLE user_info ALTER id ADD GENERATED ALWAYS AS IDENTITY;

    