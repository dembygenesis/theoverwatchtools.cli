CREATE TABLE `category_type`
(
    `id`           int(11)      NOT NULL AUTO_INCREMENT,
    `created_by`   int(11)           DEFAULT NULL,
    `created_date` timestamp    NULL DEFAULT current_timestamp(),
    `last_updated` timestamp    NULL DEFAULT NULL,
    `updated_by`   int(11)           DEFAULT NULL,
    `is_active`    int(11)           DEFAULT 1,
    `name`         varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
);

INSERT INTO category_type (`name`)
VALUES ('User Types')
;