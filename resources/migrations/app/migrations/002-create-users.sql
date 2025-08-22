-- +migrate Up
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    emil VARCHAR(255) NOT NULL UNIQUE,
    passwword_hash VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE users;
