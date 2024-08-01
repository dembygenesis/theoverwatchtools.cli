SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `user`
(
    `id`                   int(11) NOT NULL AUTO_INCREMENT,
    `firstname`            varchar(255) NOT NULL,
    `lastname`             varchar(255) NOT NULL,
    `email`                varchar(255) NOT NULL,
    `password`             varchar(255) NOT NULL,
    `organization_ref_id`  int(11) DEFAULT NULL,
    `category_type_ref_id` int(11) NOT NULL,
    `created_by`           int(11) DEFAULT NULL,
    `last_updated_by`      int(11) DEFAULT NULL,
    `created_at`           timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`      timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`            bool         NOT NULL DEFAULT TRUE,
    `reset_token`          varchar(255)          DEFAULT NULL,
    `address`              varchar(255)          DEFAULT NULL,
    `birthday`             date                  DEFAULT NULL,
    `gender`               char(1)               DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    KEY                    `user_fk_category_type` (`category_type_ref_id`),
    KEY                    `user_fk_created_by` (`created_by`),
    KEY                    `user_fk_created_last_updated_by` (`last_updated_by`),
    CONSTRAINT `user_fk_category_type` FOREIGN KEY (`category_type_ref_id`) REFERENCES `category` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `user_organization_id_fk` FOREIGN KEY (`organization_ref_id`) REFERENCES `organization` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `user_fk_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `user_fk_created_last_updated_by` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`) ON DELETE RESTRICT
);


-- Set dynamic values for category_type_ref_id and created_by
SET
@category_type_ref_super_admin = (SELECT id FROM category WHERE name = 'Super Admin');
SET
@category_type_ref_admin = (SELECT id FROM category WHERE name = 'Admin');
SET
@category_type_ref_regular = (SELECT id FROM category WHERE name = 'Regular User');
SET
@created_by1 = 1;
SET
@created_by2 = 2;
SET
@created_by3 = 3;

INSERT INTO `user`
(`firstname`, `lastname`, `email`, `password`, `category_type_ref_id`, `created_by`, `last_updated_by`, `is_active`,
 `reset_token`, `address`, `birthday`, `gender`)
VALUES ('Super Admin', 'User', 'demby@gmail.com', 'password123', @category_type_ref_super_admin, NULL, NULL, TRUE, NULL,
        '123 Main St', '1990-01-01', 'M'),
       ('Admin', 'User', 'demby@yahoo.com', 'password123', @category_type_ref_admin, @created_by1, NULL, TRUE, NULL,
        '456 Elm St', '1985-05-15', 'F'),
       ('Alice', 'User', 'demby@hotmail.com', 'password123', @category_type_ref_regular, @created_by1, NULL, TRUE, NULL,
        '789 Oak St', '1992-07-21', 'F');

SET
FOREIGN_KEY_CHECKS = 1;
