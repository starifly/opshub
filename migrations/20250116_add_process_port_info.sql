-- 添加进程和端口详情字段到hosts表
ALTER TABLE `hosts` ADD COLUMN `process_info` TEXT NULL COMMENT '进程信息JSON' AFTER `process_count`;
ALTER TABLE `hosts` ADD COLUMN `port_info` TEXT NULL COMMENT '端口信息JSON' AFTER `port_count`;
