CREATE TABLE `capture_page_sets` (
                                     `id` int(11) NOT NULL AUTO_INCREMENT,
                                     `name` VARCHAR(255) NOT NULL,
                                     `url_name` VARCHAR(255) NOT NULL,
                                     `created_by` int(11) NOT NULL,
                                     `updated_by` int(11) NOT NULL,
                                     `organization_id` INT NOT NULL,
                                     `switch_duration_in_minutes` INT NOT NULL DEFAULT 0,
                                     `created_at` TIMESTAMP NULL DEFAULT NULL,
                                     `updated_at` TIMESTAMP NULL DEFAULT NULL,
                                     `analytics_number_of_forms` INT NOT NULL DEFAULT 0 COMMENT 'Processed by worker queue',
                                     `analytics_impressions` INT NOT NULL DEFAULT 0 COMMENT 'Processed by worker queue',
                                     `analytics_submissions` INT NOT NULL DEFAULT 0 COMMENT 'Processed by worker queue',
                                     `analytics_last_updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     `deleted_at` TIMESTAMP NULL DEFAULT NULL,
                                     PRIMARY KEY (`id`),
                                     UNIQUE KEY `unique_name_organization` (`name`, `organization_id`),
                                     KEY `created_by` (`created_by`),
                                     KEY `updated_by` (`updated_by`),
                                     KEY `organization_id` (`organization_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `capture_page_sets` (
    `name`, `url_name`, `created_by`, `updated_by`, `organization_id`
) VALUES
      ('Test Page Set 1', 'test-page-set-1', 1, 1, 1),
      ('Test Page Set 2', 'test-page-set-2', 1, 1, 1),
      ('Test Page Set 3', 'test-page-set-3', 1, 1, 1);