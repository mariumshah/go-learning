-- +migrate Up
create schema if not exists playground CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;


create user if not exists 'playground_user'@'localhost'
identified by 'playground_password';

-- to give user access to any host:
create user if not exists 'playground_user'@'%'
identified by 'playground_password';

-- grant SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES
--  ON playground.* to 'playground_user'@'localhost';
-- grant SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES
-- ON playground.* to 'playground_user'@'%';
-- grant SELECT, INSERT, UPDATE, DELETE
--  on playground.* to 'playground_user'@'%';


-- +migrate Down
revoke all on playground.* from 'playground_user'@'localhost';
-- revoke all on playground.* from 'playground_user'@'&';
drop user if exists 'playground_user'@'localhost';
drop user if exists 'playground_user'@'&';
drop schema if exists playground;
