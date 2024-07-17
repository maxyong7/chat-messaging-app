ALTER TABLE user_credentials
ADD user_uuid TEXT DEFAULT gen_random_uuid (); 

ALTER TABLE user_info DROP CONSTRAINT user_info_pkey;
ALTER TABLE user_info RENAME COLUMN unique_id TO user_uuid;
ALTER TABLE user_info ADD id BIGINT GENERATED BY DEFAULT AS IDENTITY;
ALTER TABLE user_info ADD PRIMARY KEY (id);

ALTER TABLE reaction ADD CONSTRAINT idx_reaction_msg_user UNIQUE (message_uuid, user_uuid);

    