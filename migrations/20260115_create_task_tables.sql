-- 任务中心相关表

-- 任务作业表
CREATE TABLE IF NOT EXISTS `job_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '任务名称',
  `template_id` bigint unsigned DEFAULT NULL COMMENT '关联模板ID',
  `task_type` varchar(50) NOT NULL COMMENT '任务类型: manual-手动任务, ansible-ansible任务, cron-定时任务',
  `status` varchar(50) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待执行, running-执行中, success-成功, failed-失败',
  `target_hosts` text COMMENT '目标主机列表JSON',
  `parameters` json COMMENT '任务参数',
  `execute_time` datetime DEFAULT NULL COMMENT '执行时间',
  `result` json COMMENT '执行结果',
  `error_message` text COMMENT '错误信息',
  `created_by` bigint unsigned NOT NULL COMMENT '创建人ID',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_task_type` (`task_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务作业表';

-- 任务模板表
CREATE TABLE IF NOT EXISTS `job_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '模板名称',
  `code` varchar(100) NOT NULL COMMENT '模板编码',
  `description` text COMMENT '描述',
  `content` longtext NOT NULL COMMENT '模板内容',
  `variables` json COMMENT '变量定义JSON',
  `category` varchar(50) NOT NULL COMMENT '分类: script-脚本模板, ansible-ansible模板, module-模块模板',
  `platform` varchar(50) DEFAULT NULL COMMENT '适用平台: linux-脚本, windows-Windows脚本',
  `timeout` int DEFAULT 300 COMMENT '超时时间(秒)',
  `sort` int DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `created_by` bigint unsigned NOT NULL COMMENT '创建人ID',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_code` (`code`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务模板表';

-- Ansible任务表
CREATE TABLE IF NOT EXISTS `ansible_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '任务名称',
  `playbook_content` longtext COMMENT 'Playbook内容',
  `playbook_path` varchar(500) COMMENT 'Playbook文件路径',
  `inventory` text COMMENT 'Inventory配置JSON',
  `extra_vars` json COMMENT '额外变量JSON',
  `tags` varchar(500) COMMENT '标签列表,逗号分隔',
  `fork` int DEFAULT 5 COMMENT '并发数',
  `timeout` int DEFAULT 600 COMMENT '超时时间(秒)',
  `verbose` varchar(20) DEFAULT 'v' COMMENT '输出级别: v, vv, vvv',
  `status` varchar(50) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待执行, running-执行中, success-成功, failed-失败, cancelled-已取消',
  `last_run_time` datetime DEFAULT NULL COMMENT '最后执行时间',
  `last_run_result` json COMMENT '最后执行结果',
  `created_by` bigint unsigned NOT NULL COMMENT '创建人ID',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_last_run_time` (`last_run_time`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Ansible任务表';
