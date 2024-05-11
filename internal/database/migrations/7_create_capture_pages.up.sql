CREATE TABLE capture_pages (
                               id INT AUTO_INCREMENT PRIMARY KEY,
                               name VARCHAR(255) UNIQUE NOT NULL,
                               html LONGTEXT,
                               clicks INT DEFAULT 0,
                               is_control TINYINT(1) DEFAULT 0 NOT NULL,
                               capture_page_set_id INT,
                               created_by INT,
                               updated_by INT,
                               last_impression_at TIMESTAMP NULL,
                               impressions INT DEFAULT 0,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               deleted_at TIMESTAMP NULL,
                               FOREIGN KEY (capture_page_set_id) REFERENCES capture_page_sets(id) ON DELETE RESTRICT,
                               FOREIGN KEY (created_by) REFERENCES user(id) ON DELETE RESTRICT,
                               FOREIGN KEY (updated_by) REFERENCES user(id) ON DELETE RESTRICT,
                               UNIQUE (name, capture_page_set_id)
);
