SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME = 'idx_subscribe_type_show_sort');
SET @sql = IF(@index_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP INDEX `idx_subscribe_type_show_sort`',
    'SELECT ''Index idx_subscribe_type_show_sort does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME = 'idx_subscribe_type_sell_sort');
SET @sql = IF(@index_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP INDEX `idx_subscribe_type_sell_sort`',
    'SELECT ''Index idx_subscribe_type_sell_sort does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME = 'idx_subscribe_code');
SET @sql = IF(@index_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP INDEX `idx_subscribe_code`',
    'SELECT ''Index idx_subscribe_code does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'option_type');
SET @sql = IF(@column_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP COLUMN `option_type`',
    'SELECT ''Column option_type does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'code');
SET @sql = IF(@column_exists > 0,
    'ALTER TABLE `subscribe_price_option` DROP COLUMN `code`',
    'SELECT ''Column code does not exist in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
