CREATE TABLE click_trackers (
                                id INT AUTO_INCREMENT PRIMARY KEY,
                                name VARCHAR(255) NOT NULL UNIQUE,
                                url_name VARCHAR(255) UNIQUE,
                                redirect_url VARCHAR(255),
                                clicks INT DEFAULT 0,
                                unique_clicks INT DEFAULT 0,
                                created_by INT,
                                updated_by INT,
                                click_tracker_set_id INT,
                                country_id INT NULL,
                                created_at TIMESTAMP NULL,
                                updated_at TIMESTAMP NULL,
                                deleted_at TIMESTAMP NULL,
                                FOREIGN KEY (created_by) REFERENCES user(id) ON DELETE RESTRICT,
                                FOREIGN KEY (updated_by) REFERENCES user(id) ON DELETE RESTRICT,
                                FOREIGN KEY (click_tracker_set_id) REFERENCES click_tracker_sets(id) ON DELETE RESTRICT,
                                FOREIGN KEY (country_id) REFERENCES countries(id)
);