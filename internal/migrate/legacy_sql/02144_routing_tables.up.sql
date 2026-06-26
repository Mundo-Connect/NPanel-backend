CREATE TABLE IF NOT EXISTS `routing_profile` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Routing Profile ID',
    `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Stable profile code',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Profile display name',
    `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'Profile description',
    `scope_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'global' COMMENT 'user/plan/group/node/global',
    `scope_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'default' COMMENT 'Scope identifier',
    `priority` bigint NOT NULL DEFAULT 100 COMMENT 'Lower priority matches first',
    `mode` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'observe' COMMENT 'off/observe/enforce',
    `enabled` bool NOT NULL DEFAULT true COMMENT 'Profile enabled',
    `profile_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 profile object JSON',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    UNIQUE KEY `routing_profile_code_key` (`code`),
    KEY `routing_profile_priority_updated_idx` (`priority`, `updated_at`),
    KEY `routing_profile_scope_idx` (`scope_type`, `scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `routing_rule` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Routing Rule ID',
    `profile_id` bigint NOT NULL DEFAULT 0 COMMENT 'Bound routing_profile.id, 0 means default P1 profile',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Rule display name',
    `priority` bigint NOT NULL DEFAULT 100 COMMENT 'Lower priority matches first',
    `enabled` bool NOT NULL DEFAULT true COMMENT 'Rule enabled',
    `service_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Unlock service code',
    `matcher_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 matcher JSON',
    `action_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 action JSON',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    KEY `routing_rule_profile_priority_idx` (`profile_id`, `priority`, `updated_at`),
    KEY `routing_rule_service_code_idx` (`service_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `routing_dns_resolver` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'DNS Resolver ID',
    `tag` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Stable resolver tag',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Resolver display name',
    `proto` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'doh' COMMENT 'doh/dot/udp/tcp',
    `address` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Resolver address',
    `port` bigint NOT NULL DEFAULT 443 COMMENT 'Resolver port',
    `enabled` bool NOT NULL DEFAULT true COMMENT 'Resolver enabled',
    `resolver_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 DNS resolver JSON',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    UNIQUE KEY `routing_dns_resolver_tag_key` (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `routing_outbound` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Route Outbound ID',
    `tag` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Stable outbound tag',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Outbound display name',
    `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'node_group' COMMENT 'node/node_group/external',
    `region` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Region code',
    `enabled` bool NOT NULL DEFAULT true COMMENT 'Outbound enabled',
    `outbound_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 outbound JSON',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    UNIQUE KEY `routing_outbound_tag_key` (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `routing_unlock_service` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Unlock Service ID',
    `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Stable service code',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Service display name',
    `category` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Service category',
    `enabled` bool NOT NULL DEFAULT true COMMENT 'Service enabled',
    `service_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'routing_profile.v1 unlock service JSON',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    UNIQUE KEY `routing_unlock_service_code_key` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
