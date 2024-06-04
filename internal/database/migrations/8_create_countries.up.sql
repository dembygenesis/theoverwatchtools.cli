CREATE TABLE countries (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           code VARCHAR(255) UNIQUE,
                           name VARCHAR(255) UNIQUE,
                           created_at TIMESTAMP NULL,
                           updated_at TIMESTAMP NULL,
                           deleted_at TIMESTAMP NULL,
                           INDEX idx_code (code),
                           INDEX idx_name (name)
);
