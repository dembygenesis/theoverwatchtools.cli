SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `click_tracker`
(
    `id`                   int(11) NOT NULL AUTO_INCREMENT,
    `name`                 varchar(255) NOT NULL,
    `url_name`             LONGTEXT,

    `redirect_url`         INTEGER      NOT NULL DEFAULT 0,
    `clicks`               INTEGER      NOT NULL DEFAULT 0,
    `unique_clicks`        INTEGER      NOT NULL DEFAULT 0,
    `last_impression_at`   timestamp NULL DEFAULT NULL,
    `click_tracker_set_id` INTEGER      NOT NULL,

    -- Audit fields
    `created_by`           int(11) DEFAULT NULL,
    `last_updated_by`      int(11) DEFAULT NULL,
    `created_at`           timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`      timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`            bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `click_tracker_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `click_tracker_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),
    CONSTRAINT `click_tracker_click_tracker_set_id_if` FOREIGN KEY (`click_tracker_set_id`) REFERENCES `click_tracker_set` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY click_tracker_unique_name_organization (name, click_tracker_set_id)
);

SET
FOREIGN_KEY_CHECKS = 1;
