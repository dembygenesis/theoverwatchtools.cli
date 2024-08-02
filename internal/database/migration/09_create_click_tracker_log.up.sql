SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `click_tracker_log`
(
    `id`               int(11) NOT NULL AUTO_INCREMENT,
    `name`             varchar(255) NOT NULL,
    `ip_address`       LONGTEXT,
    `redirect_url`     INTEGER      NOT NULL DEFAULT 0,
    `details`          INTEGER      NOT NULL DEFAULT 0,
    `click_tracker_id` INTEGER      NOT NULL,

    -- Audit fields
    `created_by`       int(11) DEFAULT NULL,
    `last_updated_by`  int(11) DEFAULT NULL,
    `created_at`       timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`  timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`        bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `click_tracker_log_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `click_tracker_log_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),
    CONSTRAINT `click_tracker_log_click_tracker_id_fk` FOREIGN KEY (`click_tracker_id`) REFERENCES `click_tracker` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY click_tracker_log_name_click_tracker_set_id (name, click_tracker_id)
);

SET
FOREIGN_KEY_CHECKS = 1;