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


INSERT INTO capture_pages (
    name,
    html,
    clicks,
    is_control,
    capture_page_set_id,
    created_by,
    updated_by,
    last_impression_at,
    impressions
) VALUES
      ('Dummy Capture Page 1', '<p>This is dummy capture page 1.</p>', 10, 1, 1, 1, 1, '2024-05-01 10:00:00', 100),
      ('Dummy Capture Page 2', '<p>This is dummy capture page 2.</p>', 5, 0, 2, 2, 2, '2024-05-02 12:00:00', 50);
