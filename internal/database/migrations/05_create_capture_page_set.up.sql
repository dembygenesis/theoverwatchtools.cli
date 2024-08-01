SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `capture_page_set`
(
    `id`                        int(11) NOT NULL AUTO_INCREMENT,
    `name`                      varchar(255) NOT NULL,
    `url_name`                  LONGTEXT,
    `switch_duration`           INTEGER      NOT NULL,
    `organization_ref_id`       INTEGER,

    -- Updated by workers
    `analytics_number_of_forms` INTEGER      NOT NULL DEFAULT 0,
    `analytics_impressions`     INTEGER      NOT NULL DEFAULT 0,
    `analytics_submissions`     INTEGER      NOT NULL DEFAULT 0,
    `analytics_last_updated_at` INTEGER      NOT NULL DEFAULT 0,

    -- Audit fields
    `created_by`                int(11) DEFAULT NULL,
    `last_updated_by`           int(11) DEFAULT NULL,
    `created_at`                timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`           timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`                 bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `capture_page_set_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_page_set_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_page_set_organization_id_fk` FOREIGN KEY (`organization_ref_id`) REFERENCES `organization` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY capture_page_set_unique_name_organization (name, organization_ref_id)
);

SET
FOREIGN_KEY_CHECKS = 1;
