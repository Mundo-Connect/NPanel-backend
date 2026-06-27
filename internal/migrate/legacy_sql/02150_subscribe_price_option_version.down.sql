SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'version');
SET @sql = IF(@column_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP COLUMN `version`',
    'SELECT ''Column version does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
