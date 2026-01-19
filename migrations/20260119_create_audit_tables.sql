-- 创建操作日志和登录日志表
-- 执行时间: 2025-01-19

-- 1. 创建操作日志表
CREATE TABLE IF NOT EXISTS sys_operation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    -- 用户信息
    user_id INTEGER NOT NULL DEFAULT 0,
    username VARCHAR(50) NOT NULL DEFAULT '',
    real_name VARCHAR(50) NOT NULL DEFAULT '',

    -- 操作信息
    module VARCHAR(50) NOT NULL DEFAULT '',
    action VARCHAR(50) NOT NULL DEFAULT '',
    description VARCHAR(200) NOT NULL DEFAULT '',

    -- 请求信息
    method VARCHAR(10) NOT NULL DEFAULT '',
    path VARCHAR(200) NOT NULL DEFAULT '',
    params TEXT,

    -- 响应信息
    status INTEGER NOT NULL DEFAULT 0,
    error_msg TEXT,
    cost_time INTEGER NOT NULL DEFAULT 0,

    -- 环境信息
    ip VARCHAR(50) NOT NULL DEFAULT '',
    user_agent VARCHAR(500) NOT NULL DEFAULT ''
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_operation_log_user_id ON sys_operation_log(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_log_username ON sys_operation_log(username);
CREATE INDEX IF NOT EXISTS idx_operation_log_module ON sys_operation_log(module);
CREATE INDEX IF NOT EXISTS idx_operation_log_action ON sys_operation_log(action);
CREATE INDEX IF NOT EXISTS idx_operation_log_created_at ON sys_operation_log(created_at);
CREATE INDEX IF NOT EXISTS idx_operation_log_deleted_at ON sys_operation_log(deleted_at);

-- 2. 创建登录日志表
CREATE TABLE IF NOT EXISTS sys_login_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    -- 用户信息
    user_id INTEGER NOT NULL DEFAULT 0,
    username VARCHAR(50) NOT NULL DEFAULT '',
    real_name VARCHAR(50) NOT NULL DEFAULT '',

    -- 登录信息
    login_type VARCHAR(20) NOT NULL DEFAULT 'web',
    login_status VARCHAR(20) NOT NULL DEFAULT '',
    login_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logout_time DATETIME,

    -- 环境信息
    ip VARCHAR(50) NOT NULL DEFAULT '',
    location VARCHAR(100) NOT NULL DEFAULT '',
    user_agent VARCHAR(500) NOT NULL DEFAULT '',

    -- 失败原因
    fail_reason VARCHAR(200) NOT NULL DEFAULT ''
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_login_log_user_id ON sys_login_log(user_id);
CREATE INDEX IF NOT EXISTS idx_login_log_username ON sys_login_log(username);
CREATE INDEX IF NOT EXISTS idx_login_log_login_type ON sys_login_log(login_type);
CREATE INDEX IF NOT EXISTS idx_login_log_login_status ON sys_login_log(login_status);
CREATE INDEX IF NOT EXISTS idx_login_log_login_time ON sys_login_log(login_time);
CREATE INDEX IF NOT EXISTS idx_login_log_deleted_at ON sys_login_log(deleted_at);
