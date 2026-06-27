SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'version');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD COLUMN `version` int NOT NULL DEFAULT 1 COMMENT ''Optimistic Lock Version'' AFTER `sort`',
    'SELECT ''Column version already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `subscribe_price_option`
SET `version` = 1
WHERE `version` <= 0;
