-- 添加操作审计菜单
-- 注意：执行前请根据实际的role_id修改下面的值

-- 1. 插入"操作审计"父菜单
INSERT INTO sys_menu (name, code, type, parent_id, path, component, icon, sort, visible, status, created_at, updated_at)
VALUES ('操作审计', 'audit', 1, 0, '/audit', '', 'Document', 50, 1, 1, NOW(), NOW());

-- 获取刚插入的操作审计菜单ID（MySQL）
SET @audit_menu_id = LAST_INSERT_ID();

-- 2. 插入"操作日志"子菜单
INSERT INTO sys_menu (name, code, type, parent_id, path, component, icon, sort, visible, status, created_at, updated_at)
VALUES ('操作日志', 'operation-logs', 2, @audit_menu_id, '/audit/operation-logs', 'audit/OperationLogs', 'Document', 1, 1, 1, NOW(), NOW());

-- 3. 插入"登录日志"子菜单
INSERT INTO sys_menu (name, code, type, parent_id, path, component, icon, sort, visible, status, created_at, updated_at)
VALUES ('登录日志', 'login-logs', 2, @audit_menu_id, '/audit/login-logs', 'audit/LoginLogs', 'CircleCheck', 2, 1, 1, NOW(), NOW());

-- 4. 插入"数据日志"子菜单
INSERT INTO sys_menu (name, code, type, parent_id, path, component, icon, sort, visible, status, created_at, updated_at)
VALUES ('数据日志', 'data-logs', 2, @audit_menu_id, '/audit/data-logs', 'audit/DataLogs', 'DataLine', 3, 1, 1, NOW(), NOW());

-- 5. 为超级管理员角色分配这些菜单权限
-- 请将下面的1替换为实际的超级管理员role_id
INSERT INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE code IN ('audit', 'operation-logs', 'login-logs', 'data-logs')
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 查询验证
SELECT id, name, code, type, parent_id, path, component, icon, sort
FROM sys_menu
WHERE code IN ('audit', 'operation-logs', 'login-logs', 'data-logs')
ORDER BY parent_id, sort;
