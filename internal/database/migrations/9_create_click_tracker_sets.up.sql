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

INSERT INTO click_tracker_sets (name, url_name, created_by, updated_by, organization_id, analytics_number_of_links, analytics_last_updated_at, created_at, updated_at)
VALUES
    ('Tracker Set 1', 'tracker-set-1', 1, 1, 1, 10, NOW(), NOW(), NOW()),
    ('Tracker Set 2', 'tracker-set-2', 2, 2, 2, 20, NOW(), NOW(), NOW()),
    ('Tracker Set 3', 'tracker-set-3', 2, 2, 2, 20, NOW(), NOW(), NOW()),
    ('Tracker Set 4', 'tracker-set-4', 3, 3, 3, 30, NOW(), NOW(), NOW());
