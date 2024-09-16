CREATE TABLE `friend` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `remark` varchar(255) DEFAULT NULL,
    `add_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `friend_request` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `req_msg` varchar(255) DEFAULT NULL,
    `req_time` timestamp  NOT NULL,
    `handle_result`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_msg` varchar(255) DEFAULT NULL,
    `handled_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `group` (
    `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `icon` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `status`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `group_type` int(11) NOT NULL ,
    `is_verify` boolean NOT NULL ,
    `notification` varchar(255) DEFAULT NULL,
    `notification_uid` varchar(64) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `group_member` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `role_level`  tinyint COLLATE utf8mb4_unicode_ci NOT NULL ,
    `join_time` timestamp NULL DEFAULT NULL,
    `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `inviter_uid` varchar(64) DEFAULT NULL,
    `operator_uid` varchar(64) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `group_request` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `req_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
    `req_msg` varchar(255) DEFAULT NULL,
    `req_time` timestamp NULL DEFAULT NULL,
    `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `inviter_user_id` varchar(64) DEFAULT NULL,
    `handle_user_id` varchar(64) DEFAULT NULL,
    `handle_time` timestamp NULL DEFAULT NULL,
    `handle_result`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
