create database beyond_qa;
use beyond_qa;

CREATE TABLE `question`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '问题ID',
    `user_id`     bigint(20) UNSIGNED NOT NULL COMMENT '提问者',
    `title`       varchar(255)        NOT NULL COMMENT '问题标题',
    `description` text COMMENT '问题描述',
    `is_resolved` boolean             NOT NULL DEFAULT false COMMENT '是否已解决',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT = '问题';

CREATE TABLE `answer`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '回答ID',
    `question_id` bigint(20) UNSIGNED NOT NULL COMMENT '关联问题',
    `user_id`     bigint(20) UNSIGNED NOT NULL COMMENT '回答者',
    `content`     text COLLATE utf8mb4_bin COMMENT '回答内容',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `ix_question_id` (`question_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT = '回答';
