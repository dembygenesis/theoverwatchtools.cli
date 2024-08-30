SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `capture_page`
(
    `id`                   int(11) NOT NULL AUTO_INCREMENT,
    `name`                 varchar(255) NOT NULL,
    `html`                 LONGTEXT DEFAULT NULL,
    `clicks`               INTEGER NOT NULL DEFAULT 0,
    `capture_page_set_id`  INTEGER NOT NULL,
    `is_control`           TINYINT NOT NULL DEFAULT 0,
    `impressions`          int(11) DEFAULT 0,
    `last_impression_at`   timestamp NULL DEFAULT NULL,

    -- Audit fields
    `created_by`           int(11) DEFAULT NULL,
    `last_updated_by`      int(11) DEFAULT NULL,
    `created_at`           timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`      timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`            bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `capture_page_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_page_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_capture_page_set_id_fk` FOREIGN KEY (`capture_page_set_id`) REFERENCES `capture_page_set` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY capture_page_unique_name_organization (name, capture_page_set_id)
);

SET FOREIGN_KEY_CHECKS = 1;
