
-- +migrate Up
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES ON playground.* to 'playground_user'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES ON playground.* to 'playground_user'@'%';

-- +migrate Down
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON playground.* to 'playground_user'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON playground.* to 'playground_user'@'%';
