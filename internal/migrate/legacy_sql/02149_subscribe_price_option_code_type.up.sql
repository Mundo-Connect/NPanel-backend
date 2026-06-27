SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'code');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD COLUMN `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '''' COMMENT ''Stable Option Code'' AFTER `subscribe_id`',
    'SELECT ''Column code already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND COLUMN_NAME = 'option_type');
SET @sql = IF(@column_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD COLUMN `option_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT ''duration'' COMMENT ''Option Type'' AFTER `code`',
    'SELECT ''Column option_type already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `subscribe_price_option`
SET `option_type` = 'duration'
WHERE `option_type` = '';

UPDATE `subscribe_price_option`
SET `code` = CASE
    WHEN `duration_unit` = 'NoLimit' THEN CONCAT('duration_no_limit_', `id`)
    ELSE CONCAT('duration_', `duration_value`, '_', LOWER(`duration_unit`), '_', `id`)
END
WHERE `code` = '';

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME IN ('idx_subscribe_code', 'proxysubscribepriceoption_subscribe_id_code'));
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD INDEX `idx_subscribe_code` (`subscribe_id`, `code`)',
    'SELECT ''Index idx_subscribe_code already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME IN ('idx_subscribe_type_sell_sort', 'proxysubscribepriceoption_subscribe_id_option_type_sell_sort'));
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD INDEX `idx_subscribe_type_sell_sort` (`subscribe_id`, `option_type`, `sell`, `sort`)',
    'SELECT ''Index idx_subscribe_type_sell_sort already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscribe_price_option' AND INDEX_NAME IN ('idx_subscribe_type_show_sort', 'proxysubscribepriceoption_subscribe_id_option_type_show_sort'));
SET @sql = IF(@index_exists = 0,
    'ALTER TABLE `subscribe_price_option` ADD INDEX `idx_subscribe_type_show_sort` (`subscribe_id`, `option_type`, `show`, `sort`)',
    'SELECT ''Index idx_subscribe_type_show_sort already exists in subscribe_price_option table''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
