// Sliver 服务端 - 守护进程
// 处理客户端连接、会话管理、载荷生成等核心功能

use anyhow::Result;
use clap::{Parser, Subcommand};
use rust_i18n::t;
use sliver_shared::{Session, ListenerConfig, AppConfig};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tokio;
use std::time::Duration;

// 初始化 i18n
rust_i18n::i18n!("locales", fallback = "en");

/// Sliver 跨平台远程管理工具 - 服务端
#[derive(Parser, Debug)]
#[command(name = "sliver-server")]
#[command(author = "nnewkin")]
#[command(version = "0.1.0")]
#[command(about = t!("server.name").to_string(), long_about = t!("server.description").to_string())]
struct Args {
    /// 监听地址
    #[arg(short = 'H', long, default_value = "0.0.0.0")]
    host: String,

    /// 监听端口
    #[arg(short, long, default_value_t = 31337)]
    port: u16,

    /// 语言设置 (zh-CN, en)
    #[arg(long, default_value = "zh-CN")]
    language: String,

    /// 日志级别
    #[arg(short, long, default_value = "info")]
    log_level: String,

    /// 子命令
    #[command(subcommand)]
    command: Option<Commands>,
}

/// 服务端子命令
#[derive(Subcommand, Debug)]
enum Commands {
    /// 启动服务器
    Start,

    /// 停止服务器
    Stop,

    /// 重启服务器
    Restart,

    /// 状态
    Status,

    /// 生成载荷
    Generate {
        /// 载荷类型
        #[arg(short, long)]
        payload_type: String,

        /// 目标操作系统
        #[arg(short, long)]
        os: String,

        /// 目标架构
        #[arg(short, long)]
        arch: String,
    },

    /// 监听器管理
    Listener {
        /// 子命令
        #[command(subcommand)]
        listener_command: ListenerCommands,
    },
}

/// 监听器管理命令
#[derive(Subcommand, Debug)]
enum ListenerCommands {
    /// 列出监听器
    List,

    /// 创建监听器
    Create {
        /// 监听器类型
        #[arg(short, long)]
        listener_type: String,

        /// 监听地址
        #[arg(short, long)]
        host: String,

        /// 监听端口
        #[arg(short, long)]
        port: u16,
    },

    /// 删除监听器
    Remove {
        /// 监听器名称
        name: String,
    },
}

/// 服务端状态
#[derive(Debug, Clone)]
enum ServerState {
    Stopped,
    Starting,
    Running,
    Stopping,
}

/// 服务端核心结构
struct Server {
    config: AppConfig,
    state: Arc<Mutex<ServerState>>,
    sessions: Arc<Mutex<HashMap<String, Session>>>,
    listeners: Arc<Mutex<HashMap<String, ListenerConfig>>>,
}

impl Server {
    /// 创建新的服务端实例
    fn new(config: AppConfig) -> Self {
        Self {
            config,
            state: Arc::new(Mutex::new(ServerState::Stopped)),
            sessions: Arc::new(Mutex::new(HashMap::new())),
            listeners: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    /// 启动服务器
    async fn start(&self) -> Result<()> {
        tracing::info!("{}", t!("server.starting"));

        // 更新状态
        {
            let mut state = self.state.lock().unwrap();
            *state = ServerState::Starting;
        }

        // 这里应该是启动网络监听的实际逻辑
        // 暂时模拟启动过程
        
        tokio::time::sleep(Duration::from_millis(500)).await;

        // 更新状态为运行中
        {
            let mut state = self.state.lock().unwrap();
            *state = ServerState::Running;
        }

        tracing::info!(
            "{} {}:{}", 
            t!("server.listening_on"), 
            self.config.server_url, 
            self.config.server_port
        );
        
        tracing::info!("{}", t!("server.started"));

        Ok(())
    }

    /// 停止服务器
    async fn stop(&self) -> Result<()> {
        tracing::info!("{}", t!("server.stopping"));

        // 更新状态
        {
            let mut state = self.state.lock().unwrap();
            *state = ServerState::Stopping;
        }

        // 这里应该是停止网络监听的实际逻辑
        // 暂时模拟停止过程
        
        tokio::time::sleep(Duration::from_millis(500)).await;

        // 清理会话
        {
            let mut sessions = self.sessions.lock().unwrap();
            sessions.clear();
        }

        // 更新状态为已停止
        {
            let mut state = self.state.lock().unwrap();
            *state = ServerState::Stopped;
        }

        tracing::info!("{}", t!("server.stopped"));

        Ok(())
    }

    /// 重启服务器
    async fn restart(&self) -> Result<()> {
        self.stop().await?;
        tokio::time::sleep(Duration::from_millis(1000)).await;
        self.start().await?;
        Ok(())
    }

    /// 获取服务器状态
    fn get_status(&self) -> String {
        let state = self.state.lock().unwrap();
        match *state {
            ServerState::Stopped => t!("status.stopped").to_string(),
            ServerState::Starting => t!("status.connecting").to_string(),
            ServerState::Running => t!("status.running").to_string(),
            ServerState::Stopping => t!("status.pending").to_string(),
        }
    }

    /// 添加会话（暂未使用）
    #[allow(dead_code)]
    fn add_session(&self, session: Session) -> Result<()> {
        let mut sessions = self.sessions.lock().unwrap();
        sessions.insert(session.id.clone(), session);
        Ok(())
    }
    
    /// 移除会话（暂未使用）
    #[allow(dead_code)]
    fn remove_session(&self, session_id: &str) -> Result<()> {
        let mut sessions = self.sessions.lock().unwrap();
        if sessions.remove(session_id).is_some() {
            Ok(())
        } else {
            Err(anyhow::anyhow!("Session not found"))
        }
    }

    /// 获取所有会话
    fn get_sessions(&self) -> Vec<Session> {
        let sessions = self.sessions.lock().unwrap();
        sessions.values().cloned().collect()
    }

    /// 添加监听器
    fn add_listener(&self, listener: ListenerConfig) -> Result<()> {
        let mut listeners = self.listeners.lock().unwrap();
        listeners.insert(listener.name.clone(), listener);
        Ok(())
    }

    /// 移除监听器
    fn remove_listener(&self, name: &str) -> Result<()> {
        let mut listeners = self.listeners.lock().unwrap();
        if listeners.remove(name).is_some() {
            Ok(())
        } else {
            Err(anyhow::anyhow!("Listener not found"))
        }
    }

    /// 获取所有监听器
    fn get_listeners(&self) -> Vec<ListenerConfig> {
        let listeners = self.listeners.lock().unwrap();
        listeners.values().cloned().collect()
    }

    /// 运行服务器主循环
    async fn run(&self) -> Result<()> {
        tracing::info!("{}", t!("server.started"));
        
        // 这里应该是主事件循环
        // 暂时模拟运行
        
        tokio::signal::ctrl_c().await?;
        
        tracing::info!("Received shutdown signal");
        self.stop().await?;
        
        Ok(())
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    // 解析命令行参数
    let args = Args::parse();

    // 设置语言
    let locale = match args.language.as_str() {
        "zh-CN" => "zh-CN",
        "en" => "en",
        _ => "zh-CN",
    };
    rust_i18n::set_locale(locale);

    // 初始化日志
    let log_level = match args.log_level.as_str() {
        "trace" => tracing::Level::TRACE,
        "debug" => tracing::Level::DEBUG,
        "info" => tracing::Level::INFO,
        "warn" => tracing::Level::WARN,
        "error" => tracing::Level::ERROR,
        _ => tracing::Level::INFO,
    };
    
    tracing_subscriber::fmt()
        .with_max_level(log_level)
        .init();

    // 创建服务端配置
    let config = AppConfig {
        server_url: args.host.clone(),
        server_port: args.port,
        language: args.language.clone(),
        log_level: args.log_level.clone(),
    };

    // 创建服务端实例
    let server = Server::new(config);

    // 处理子命令
    match args.command {
        Some(Commands::Start) => {
            server.start().await?;
            
            // 启动服务器后进入主循环
            server.run().await?;
        }
        
        Some(Commands::Stop) => {
            server.stop().await?;
        }
        
        Some(Commands::Restart) => {
            server.restart().await?;
        }
        
        Some(Commands::Status) => {
            let status = server.get_status();
            let sessions = server.get_sessions();
            let listeners = server.get_listeners();
            
            println!("{}: {}", t!("status.running"), status);
            println!("{}: {}", t!("session.active_sessions"), sessions.len());
            println!("{}: {}", t!("listener.active_listeners"), listeners.len());
        }
        
        Some(Commands::Generate { payload_type, os, arch }) => {
            println!("{}: {} {} ({})", 
                t!("payload.generate"), 
                payload_type, 
                os, 
                arch
            );
            // 实际载荷生成逻辑
            println!("{}", t!("payload.generated"));
        }
        
        Some(Commands::Listener { listener_command }) => {
            match listener_command {
                ListenerCommands::List => {
                    let listeners = server.get_listeners();
                    println!("{}", t!("listener.list_header"));
                    for listener in listeners {
                        println!("  {} - {} ({})", 
                            listener.name, 
                            listener.url, 
                            listener.listener_type
                        );
                    }
                }
                
                ListenerCommands::Create { listener_type, host, port } => {
                    let listener = ListenerConfig {
                        name: format!("{}_{}_{}", listener_type, host, port),
                        listener_type: listener_type.clone(),
                        url: format!("{}:{}", host, port),
                        host: host.clone(),
                        port,
                        protocol: listener_type.clone(),
                        is_active: true,
                    };
                    server.add_listener(listener)?;
                    println!("{} {}:{} ({})", 
                        t!("listener.create"), 
                        host, 
                        port, 
                        listener_type
                    );
                }
                
                ListenerCommands::Remove { name } => {
                    server.remove_listener(&name)?;
                    println!("{}: {}", t!("listener.removed"), name);
                }
            }
        }
        
        None => {
            // 没有子命令时，默认启动服务器
            server.start().await?;
            server.run().await?;
        }
    }

    Ok(())
}
