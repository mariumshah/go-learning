
-- +migrate Up

ALTER TABLE users
    RENAME COLUMN passwword_hash TO password_hash;

CREATE TABLE authors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_author_name (name)  -- Case-sensitive uniqueness
);

-- Books table (shared catalog)
CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(512) NOT NULL,
    author_id INT NOT NULL,  -- Normalized reference
    isbn VARCHAR(64) DEFAULT NULL,
    publication_year SMALLINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    total_libraries INT DEFAULT 0,  -- Counter cache
    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE RESTRICT,
    UNIQUE KEY unique_book_identity (title, author_id, isbn(17)),  -- Smart uniqueness
    FULLTEXT INDEX ft_book_title (title)  -- Title-only fulltext
) ENGINE=InnoDB;

-- Junction table for user libraries
CREATE TABLE user_books (
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, book_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    INDEX idx_added_at (added_at)
) ENGINE=InnoDB;

-- +migrate Down
DROP TABLE user_books;
DROP TABLE books;
DROP TABLE authors;
ALTER TABLE users RENAME COLUMN passwword_hash TO password_hash;
