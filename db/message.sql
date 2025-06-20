create database beyond_message;
use beyond_message;

CREATE TABLE `message`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `type`        tinyint(4)          NOT NULL DEFAULT 0 COMMENT '类型：0系统，1点赞，2评论…',
    `biz_id`      varchar(64)         NOT NULL DEFAULT '' COMMENT '业务ID（如文章、问题等）',
    `target_id`   bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '业务目标ID',
    `receiver_id` bigint(20) UNSIGNED NOT NULL COMMENT '接收者',
    `title`       varchar(128)        NOT NULL DEFAULT '' COMMENT '标题',
    `content`     text COLLATE utf8mb4_bin COMMENT '正文内容',
    `is_read`     boolean             NOT NULL DEFAULT false COMMENT '是否已读',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `ix_receiver_id` (`receiver_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT = '系统消息';
