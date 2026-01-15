-- 添加主机采集信息字段，用于存储从主机采集的系统信息
-- 包括CPU、内存、磁盘、进程、端口等信息

-- 添加 CPU 信息字段
ALTER TABLE `hosts`
ADD COLUMN `cpu_info` TEXT NULL COMMENT 'CPU信息JSON' AFTER `arch`;

ALTER TABLE `hosts`
ADD COLUMN `cpu_cores` INT DEFAULT 0 COMMENT 'CPU核心数' AFTER `cpu_info`;

ALTER TABLE `hosts`
ADD COLUMN `cpu_usage` FLOAT DEFAULT 0 COMMENT 'CPU使用率' AFTER `cpu_cores`;

-- 添加内存信息字段
ALTER TABLE `hosts`
ADD COLUMN `memory_total` BIGINT DEFAULT 0 COMMENT '内存总容量(字节)' AFTER `cpu_usage`;

ALTER TABLE `hosts`
ADD COLUMN `memory_used` BIGINT DEFAULT 0 COMMENT '已用内存(字节)' AFTER `memory_total`;

ALTER TABLE `hosts`
ADD COLUMN `memory_usage` FLOAT DEFAULT 0 COMMENT '内存使用率' AFTER `memory_used`;

-- 添加磁盘信息字段
ALTER TABLE `hosts`
ADD COLUMN `disk_total` BIGINT DEFAULT 0 COMMENT '磁盘总容量(字节)' AFTER `memory_usage`;

ALTER TABLE `hosts`
ADD COLUMN `disk_used` BIGINT DEFAULT 0 COMMENT '已用磁盘(字节)' AFTER `disk_total`;

ALTER TABLE `hosts`
ADD COLUMN `disk_usage` FLOAT DEFAULT 0 COMMENT '磁盘使用率' AFTER `disk_used`;

-- 添加进程和端口信息字段
ALTER TABLE `hosts`
ADD COLUMN `process_count` INT DEFAULT 0 COMMENT '进程数量' AFTER `disk_usage`;

ALTER TABLE `hosts`
ADD COLUMN `port_count` INT DEFAULT 0 COMMENT '端口数量' AFTER `process_count`;

-- 添加其他系统信息字段
ALTER TABLE `hosts`
ADD COLUMN `uptime` VARCHAR(100) NULL COMMENT '运行时间' AFTER `port_count`;

ALTER TABLE `hosts`
ADD COLUMN `hostname` VARCHAR(100) NULL COMMENT '主机名' AFTER `uptime`;

-- 为已有数据初始化默认值
UPDATE `hosts` SET `cpu_cores` = 0, `cpu_usage` = 0, `memory_total` = 0, `memory_used` = 0, `memory_usage` = 0, `disk_total` = 0, `disk_used` = 0, `disk_usage` = 0, `process_count` = 0, `port_count` = 0 WHERE `cpu_cores` IS NULL;
