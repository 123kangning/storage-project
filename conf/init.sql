create database if not exists file;

CREATE USER if not exists canal IDENTIFIED BY 'canal';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';

CREATE USER if not exists file IDENTIFIED BY 'file';
GRANT ALL PRIVILEGES ON file.* TO 'file'@'%';
-- GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
FLUSH PRIVILEGES;

use file;
CREATE TABLE IF NOT EXISTS `file`
(
    `id`        BIGINT       NOT NULL auto_increment,
    `name`      VARCHAR(255) NOT NULL,
    `source`    BIGINT       NOT NULL COMMENT '文件上传者的用户ID',
    `size`      INTEGER      NOT NULL,
    `hash`      CHAR(64)     NOT NULL,
    `is_delete` BOOL         NOT NULL COMMENT '删除之后该字段为1',
    `update_at` DATETIME     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 创建 user 表
CREATE TABLE IF NOT EXISTS `user`
(
    `id`       BIGINT       NOT NULL auto_increment,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 创建会话表
CREATE TABLE IF NOT EXISTS `sessions` (
    `id` VARCHAR(255) PRIMARY KEY,
    `user_id` INT NOT NULL,
    `expires_at` TIMESTAMP NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 在 sessions 表为 user_id 添加索引
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

