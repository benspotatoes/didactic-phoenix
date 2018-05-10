\connect historislack;

CREATE SCHEMA organization;
ALTER SCHEMA organization OWNER TO historislack;

CREATE TABLE organization.messages (
  id bigserial primary key,
  channel_id varchar(32),
  channel_name varchar(256),
  user_id varchar(32),
  user_name varchar(256),
  message_text text,
  message_timestamp varchar(256),
  created_at timestamp DEFAULT current_timestamp
);
ALTER TABLE organization.messages OWNER TO historislack;
