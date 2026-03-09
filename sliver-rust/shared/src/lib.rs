// Sliver 共享库
// 包含客户端和服务端共享的数据结构、协议定义等

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// 会话信息
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Session {
    pub id: String,
    pub name: String,
    pub hostname: String,
    pub username: String,
    pub transport: String,
    pub remote_address: String,
    pub os: String,
    pub arch: String,
    pub pid: u32,
    pub is_active: bool,
    pub last_seen: String,
}

/// 监听器配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ListenerConfig {
    pub name: String,
    pub listener_type: String,
    pub url: String,
    pub host: String,
    pub port: u16,
    pub protocol: String,
    pub is_active: bool,
}

/// 载荷配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadConfig {
    pub name: String,
    pub payload_type: String,
    pub os: String,
    pub arch: String,
    pub format: String,
    pub c2_transport: String,
}

/// RPC请求错误类型
#[derive(Debug, thiserror::Error)]
pub enum RpcError {
    #[error("网络连接失败: {0}")]
    ConnectionError(String),
    
    #[error("认证失败: {0}")]
    AuthenticationError(String),
    
    #[error("操作执行失败: {0}")]
    ExecutionError(String),
    
    #[error("无效参数: {0}")]
    InvalidParameter(String),
}

/// RPC响应
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RpcResponse<T> {
    pub success: bool,
    pub data: Option<T>,
    pub error: Option<String>,
}

impl<T> RpcResponse<T> {
    pub fn success(data: T) -> Self {
        Self {
            success: true,
            data: Some(data),
            error: None,
        }
    }
    
    pub fn error(error: String) -> Self {
        Self {
            success: false,
            data: None,
            error: Some(error),
        }
    }
}

/// 应用程序配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub server_url: String,
    pub server_port: u16,
    pub language: String,
    pub log_level: String,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            server_url: "localhost".to_string(),
            server_port: 31337,
            language: "zh-CN".to_string(),
            log_level: "INFO".to_string(),
        }
    }
}

/// 命令执行结果
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CommandResult {
    pub command: String,
    pub exit_code: i32,
    pub stdout: String,
    pub stderr: String,
    pub execution_time: u64,
}

/// 任务状态
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct JobStatus {
    pub id: String,
    pub name: String,
    pub description: String,
    pub status: String,
    pub created_at: String,
    pub updated_at: String,
}
