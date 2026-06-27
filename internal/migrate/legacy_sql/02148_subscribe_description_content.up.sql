SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND COLUMN_NAME = 'short_description');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe` ADD COLUMN `short_description` longtext NULL COMMENT ''Short Subscribe Description'' AFTER `description`',
    'SELECT ''Column short_description already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND COLUMN_NAME = 'features');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe` ADD COLUMN `features` longtext NULL COMMENT ''Subscribe Feature List JSON'' AFTER `short_description`',
    'SELECT ''Column features already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND COLUMN_NAME = 'detail_format');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe` ADD COLUMN `detail_format` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT ''markdown'' COMMENT ''Subscribe Detail Format'' AFTER `features`',
    'SELECT ''Column detail_format already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe' AND COLUMN_NAME = 'detail_content');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe` ADD COLUMN `detail_content` longtext NULL COMMENT ''Subscribe Detail Content'' AFTER `detail_format`',
    'SELECT ''Column detail_content already exists in subscribe table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
