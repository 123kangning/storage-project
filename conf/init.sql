create database if not exists file;

CREATE USER canal IDENTIFIED BY 'canal';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON file.* TO 'canal'@'%';

CREATE USER file IDENTIFIED BY 'file';
GRANT ALL PRIVILEGES ON file.* TO 'file'@'%';
-- GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
FLUSH PRIVILEGES;

use file;
CREATE TABLE IF NOT EXISTS `file`
(
    `id`        BIGINT       NOT NULL auto_increment,
    `name`      VARCHAR(255) NOT NULL,
    `size`      INTEGER      NOT NULL,
    `hash`      CHAR(64)     NOT NULL,
    `is_delete` BOOL         NOT NULL COMMENT '删除之后该字段为1',
    `update_at` DATETIME     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;
