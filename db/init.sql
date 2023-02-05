CREATE DATABASE IF NOT EXISTS ydktest  DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

use ydktest;

DROP TABLE if exists `config`;
CREATE TABLE `config` (
                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
                          `key` varchar(32) NOT NULL COMMENT '配置key',
                          `value` varchar(128) NOT NULL COMMENT '配置值',
                          `sys` varchar(32) NOT NULL DEFAULT '' COMMENT '配置项所属系统',
                          `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '描述',
                          `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `uni_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `config` (`key`, `value`,`sys`, `desc`)
VALUES
    ('station_task','1','pms','搬箱定时任务开关'),
    ('track_threshold','0.8','pms','轨道可容纳料箱阈值'),
    ('back_after_pick','1','pms','拣选完是否回箱');

INSERT INTO `config` (`key`, `value`,`sys`, `desc`) VALUES('max_intervals','4','pms','最迟出箱时间(小时)');
INSERT INTO `config` (`key`, `value`,`sys`, `desc`) VALUES('invert_box_task','1','pms','闲时倒箱定时任务开关');