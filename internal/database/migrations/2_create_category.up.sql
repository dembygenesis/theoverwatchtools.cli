CREATE TABLE `category`
(
    `id`                   int(11) NOT NULL AUTO_INCREMENT,
    `created_by`           int(11) DEFAULT NULL,
    `created_date`         timestamp NULL DEFAULT current_timestamp (),
    `last_updated`         timestamp NULL DEFAULT NULL,
    `updated_by`           int(11) DEFAULT NULL,
    `is_active`            int(11) DEFAULT 1 NOT NULL,
    `category_type_ref_id` int(11) NOT NULL,
    `name`                 varchar(255) NOT NULL,
    CONSTRAINT `category_type_ref_id_fk` FOREIGN KEY (`category_type_ref_id`) REFERENCES `category_type` (`id`),
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`, `category_type_ref_id`)
);

SET @user_types_ref_id = (SELECT id FROM `category_type` WHERE `name` = 'User Types');

INSERT INTO `category` (`name`, `category_type_ref_id`)
VALUES ('Super Admin', @user_types_ref_id),
       ('Admin', @user_types_ref_id),
       ('Regular User', @user_types_ref_id);