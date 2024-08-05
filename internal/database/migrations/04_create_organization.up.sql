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

-- the only way to make Test_UpdateOrganization_Fail work is to run this lines below with the go run ./cmd/migrate
-- INSERT INTO `organization` (`name`, `created_by`, `last_updated_by`, `created_at`, `last_updated_at`, `is_active`)
-- VALUES
--     ('Organization A', 1, 2, '2023-01-01 00:00:00', '2023-01-01 00:00:00', TRUE),
--     ('Organization B', 2, 3, '2023-02-01 00:00:00', '2023-02-01 00:00:00', TRUE),
--     ('Organization C', 3, 4, '2023-03-01 00:00:00', '2023-03-01 00:00:00', TRUE);