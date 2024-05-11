CREATE TABLE `organization`
(
    `id` INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `created_by` INT(11) DEFAULT NULL,
    `created_date` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP(),
    `last_updated` TIMESTAMP DEFAULT NULL,
    `updated_by` INT(11) DEFAULT NULL,
    `is_active` INT(11) DEFAULT 1 NOT NULL,
    `organization_type_ref_id` INT(11) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    CONSTRAINT `organization_type_ref_id_fk` FOREIGN KEY (`organization_type_ref_id`) REFERENCES `organization_type` (`id`),
    UNIQUE KEY `name` (`name`, `organization_type_ref_id`)
);

INSERT INTO `organization` (`name`, `organization_type_ref_id`)
VALUES
    ('TEST DATA 1', 1),
    ('TEST DATA 2', 1);

