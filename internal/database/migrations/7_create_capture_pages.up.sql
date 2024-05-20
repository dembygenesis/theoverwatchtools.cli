CREATE TABLE `capture_pages` (
                                 `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                 `name` VARCHAR(255) UNIQUE NOT NULL,
                                 `html` LONGTEXT NOT NULL,
                                 `clicks` INT DEFAULT 0 NOT NULL,
                                 `is_control` BOOLEAN DEFAULT 0 NOT NULL,
                                 `capture_page_set_id` BIGINT UNSIGNED NOT NULL, -- Matching the data type with `capture_page_sets.id`
                                 `created_by` INT UNSIGNED NOT NULL,
                                 `updated_by` INT UNSIGNED NOT NULL,
                                 `last_impression_at` TIMESTAMP NULL,
                                 `impressions` INT DEFAULT 0 NOT NULL,
                                 `created_at` TIMESTAMP NULL,
                                 `updated_at` TIMESTAMP NULL,
                                 `deleted_at` TIMESTAMP NULL,
                                 UNIQUE KEY `unique_name_capture_page_set` (`name`, `capture_page_set_id`),
                                 FOREIGN KEY (`capture_page_set_id`) REFERENCES `capture_page_sets` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
