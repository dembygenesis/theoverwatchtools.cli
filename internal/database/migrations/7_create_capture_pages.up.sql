CREATE TABLE `capture_pages` (
                                 `id` int(11) AUTO_INCREMENT PRIMARY KEY,
                                 `name` VARCHAR(255) UNIQUE NOT NULL,
                                 `html` LONGTEXT NOT NULL,
                                 `clicks` int(11) DEFAULT 0 NOT NULL,
                                 `is_control` int(11) DEFAULT 1 NOT NULL,
                                 `capture_page_set_id` INT NOT NULL, -- Matching the data type with `capture_page_sets.id`
                                 `created_by` int(11) NOT NULL,
                                 `updated_by` int(11) NOT NULL,
                                 `last_impression_at` TIMESTAMP NULL,
                                 `impressions` int(11) DEFAULT 0 NOT NULL,
                                 `created_at` TIMESTAMP NULL,
                                 `updated_at` TIMESTAMP NULL,
                                 `deleted_at` TIMESTAMP NULL,
                                 UNIQUE KEY `unique_name_capture_page_set` (`name`, `capture_page_set_id`),
                                 FOREIGN KEY (`capture_page_set_id`) REFERENCES `capture_page_sets` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `capture_pages` (
    `name`, `html`, `is_control`, `capture_page_set_id`, `created_by`, `updated_by`
) VALUES
      ('Example Page 1', '<html><body><h1>Example 1</h1></body></html>', 1, 1, 1, 1),
      ('Example Page 2', '<html><body><h1>Example 2</h1></body></html>', 0, 1, 1, 1);