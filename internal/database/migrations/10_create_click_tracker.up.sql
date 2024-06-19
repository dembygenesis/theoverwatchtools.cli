-- Create Click Trackers table
CREATE TABLE click_trackers (
                                id INT AUTO_INCREMENT PRIMARY KEY,
                                name VARCHAR(255) NOT NULL UNIQUE,
                                url_name VARCHAR(255) UNIQUE,
                                redirect_url VARCHAR(255),
                                clicks INT DEFAULT 1,
                                unique_clicks INT DEFAULT 0,
                                created_by INT NOT NULL,
                                updated_by INT NOT NULL,
                                click_tracker_set_id INT NOT NULL,
                                country_id INT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                deleted_at TIMESTAMP NULL,
                                FOREIGN KEY (created_by) REFERENCES user(id) ON DELETE RESTRICT,
                                FOREIGN KEY (updated_by) REFERENCES user(id) ON DELETE RESTRICT,
                                FOREIGN KEY (click_tracker_set_id) REFERENCES click_tracker_sets(id) ON DELETE RESTRICT,
                                FOREIGN KEY (country_id) REFERENCES countries(id)
);

-- Insert 3 entries into Click Trackers table
INSERT INTO click_trackers (name, url_name, redirect_url, created_by, updated_by, click_tracker_set_id)
VALUES
    ('Tracker 1', 'tracker-1-url', 'http://example.com/redirect1', 1, 1, 1),
    ('Tracker 2', 'tracker-2-url', 'http://example.com/redirect2', 2, 2, 2),
    ('Tracker 3', 'tracker-3-url', 'http://example.com/redirect3', 2, 2, 2),
    ('Tracker 4', 'tracker-4-url', 'http://example.com/redirect4', 3, 3, 3);
