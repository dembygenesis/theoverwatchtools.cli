create table `organization_type`
(
    `id`           int(11)      NOT NULL AUTO_INCREMENT,
    `created_by`   int(11)          DEFAULT NULL,
    `created_date` timestamp    NULL DEFAULT current_timestamp(),
    `last_updated` timestamp    NULL DEFAULT NULL,
    `updated_by`   int(11)          DEFAULT NULL,
    `is_active`    int(11)          DEFAULT 1,
    `name`         varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
);

INSERT INTO organization_type (`name`)
VALUES ('Organization Type 1'),
       ('Organization Type 2');

