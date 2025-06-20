create database beyond_chat;
use beyond_chat;

CREATE TABLE `chat_session`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '会话ID',
    `user1_id`    bigint(20) UNSIGNED NOT NULL COMMENT '用户1',
    `user2_id`    bigint(20) UNSIGNED NOT NULL COMMENT '用户2',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_pair` (`user1_id`, `user2_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT = '单聊会话';

CREATE TABLE `chat_message`
(
    `id`          bigint(20) UNSIGNED      NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `session_id`  bigint(20) UNSIGNED      NOT NULL COMMENT '所属会话',
    `sender_id`   bigint(20) UNSIGNED      NOT NULL COMMENT '发送者ID',
    `receiver_id` bigint(20) UNSIGNED      NOT NULL COMMENT '接收者ID',
    `content`     text COLLATE utf8mb4_bin NOT NULL COMMENT '内容',
    `msg_type`    tinyint(4)               NOT NULL DEFAULT 0 COMMENT '消息类型：0文字，1图片',
    `is_read`     boolean                  NOT NULL DEFAULT false COMMENT '是否已读',
    `send_time`   timestamp                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    PRIMARY KEY (`id`),
    KEY `ix_session_id` (`session_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT = '聊天消息';
