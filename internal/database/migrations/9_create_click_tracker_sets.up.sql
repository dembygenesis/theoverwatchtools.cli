CREATE TABLE click_tracker_sets (
                                    id INT AUTO_INCREMENT PRIMARY KEY,
                                    name VARCHAR(255) NOT NULL,
                                    url_name VARCHAR(255) NOT NULL,
                                    created_by INT,
                                    updated_by INT,
                                    organization_id INT,
                                    analytics_number_of_links INT DEFAULT 0 COMMENT 'Processed by worker queue',
                                    analytics_last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    created_at TIMESTAMP NULL,
                                    updated_at TIMESTAMP NULL,
                                    deleted_at TIMESTAMP NULL,
                                    UNIQUE (name, organization_id),
                                    FOREIGN KEY (created_by) REFERENCES user(id) ON DELETE RESTRICT,
                                    FOREIGN KEY (updated_by) REFERENCES user(id) ON DELETE RESTRICT,
                                    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE RESTRICT
);
