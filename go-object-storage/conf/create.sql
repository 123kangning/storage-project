CREATE TABLE IF NOT EXISTS `file` (
    `id` INTEGER NOT NULL auto_increment,
    `name` VARCHAR(255) NOT NULL,
    `size` INTEGER NOT NULL,
    `hash` CHAR(64) NOT NULL,
    `is_delete` BOOL NOT NULL COMMENT '删除之后该字段为1',
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;