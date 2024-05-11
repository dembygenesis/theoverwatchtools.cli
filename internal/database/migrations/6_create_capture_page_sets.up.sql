CREATE TABLE capture_page_sets (
                                   id INT AUTO_INCREMENT PRIMARY KEY,
                                   name VARCHAR(255) NOT NULL,
                                   url_name VARCHAR(255) NOT NULL,
                                   created_by INT,
                                   updated_by INT,
                                   organization_id INT,
                                   switch_duration_in_minutes INT DEFAULT 0 NOT NULL,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                   analytics_number_of_forms INT DEFAULT 0 COMMENT 'Processed by worker queue',
                                   analytics_impressions INT DEFAULT 0 COMMENT 'Processed by worker queue',
                                   analytics_submissions INT DEFAULT 0 COMMENT 'Processed by worker queue',
                                   analytics_last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                   UNIQUE (name, organization_id),
                                   FOREIGN KEY (created_by) REFERENCES user(id) ON DELETE RESTRICT,
                                   FOREIGN KEY (updated_by) REFERENCES user(id) ON DELETE RESTRICT,
                                   FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE RESTRICT
);
