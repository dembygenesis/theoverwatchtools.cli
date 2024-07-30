SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `organization`
(
    `id`                   int(11) NOT NULL AUTO_INCREMENT,
    `name`                 varchar(255) NOT NULL,

    -- Audit fields
    `created_by`           int(11) DEFAULT NULL,
    `last_updated_by`      int(11) DEFAULT NULL,
    `created_at`           timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`      timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`            bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `organization_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
);

SET FOREIGN_KEY_CHECKS = 1;
