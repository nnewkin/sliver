// Sliver 客户端 - CLI 命令行界面
// 包含所有用户交互逻辑，使用 rust-i18n 实现多语言支持

use anyhow::Result;
use clap::{Parser, Subcommand};
use crossterm::{
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, Clear, ClearType},
};
use rust_i18n::{t, Locale};
use sliver_shared::{Session, ListenerConfig, PayloadConfig, AppConfig};
use std::io::stdout;

// 初始化 i18n，设置语言资源路径
rust_i18n::i18n!("locales", fallback = "en");

/// Sliver 跨平台远程管理工具 - 客户端
#[derive(Parser, Debug)]
#[command(name = "sliver-client")]
#[command(author = "nnewkin")]
#[command(version = "0.1.0")]
#[command(about = t!("client.name").to_string(), long_about = t!("client.description").to_string())]
struct Args {
    /// 服务器地址
    #[arg(short, long, default_value = "localhost")]
    server: String,

    /// 服务器端口
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

/// 客户端子命令
#[derive(Subcommand, Debug)]
enum Commands {
    /// 连接到服务器
    Connect {
        /// 服务器地址
        #[arg(short, long)]
        host: Option<String>,

        /// 服务器端口
        #[arg(short, long)]
        port: Option<u16>,
    },

    /// 会话管理
    Session {
        /// 子命令
        #[command(subcommand)]
        session_command: SessionCommands,
    },

    /// 监听器管理
    Listener {
        /// 子命令
        #[command(subcommand)]
        listener_command: ListenerCommands,
    },

    /// 载荷生成
    Generate {
        /// 子命令
        #[command(subcommand)]
        generate_command: GenerateCommands,
    },

    /// 信标管理
    Beacon {
        /// 子命令
        #[command(subcommand)]
        beacon_command: BeaconCommands,
    },
}

/// 会话相关命令
#[derive(Subcommand, Debug)]
enum SessionCommands {
    /// 列出所有会话
    List,
    
    /// 连接到指定会话
    Use {
        /// 会话ID或名称
        session_id: String,
    },
    
    /// 终止会话
    Kill {
        /// 会话ID或名称
        session_id: String,
    },
    
    /// 查看会话信息
    Info {
        /// 会话ID或名称
        session_id: String,
    },
}

/// 监听器相关命令
#[derive(Subcommand, Debug)]
enum ListenerCommands {
    /// 列出所有监听器
    List,
    
    /// 创建新监听器
    Create {
        /// 监听器类型 (mtls, http, https, dns, wireguard)
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

/// 载荷生成相关命令
#[derive(Subcommand, Debug)]
enum GenerateCommands {
    /// 生成载荷
    Payload {
        /// 载荷类型 (dll, exe, shellcode, raw, service)
        #[arg(short, long)]
        payload_type: String,
        
        /// 目标操作系统
        #[arg(short, long)]
        os: String,
        
        /// 目标架构
        #[arg(short, long)]
        arch: String,
        
        /// C2传输方式
        #[arg(short, long)]
        c2_transport: String,
    },
    
    /// 列出所有载荷
    List,
}

/// 信标相关命令
#[derive(Subcommand, Debug)]
enum BeaconCommands {
    /// 列出所有信标
    List,
    
    /// 设置签到间隔
    Interval {
        /// 间隔时间（秒）
        seconds: u32,
    },
    
    /// 设置抖动百分比
    Jitter {
        /// 抖动百分比 (0-100)
        percent: u32,
    },
}

/// 客户端主结构
struct Client {
    config: AppConfig,
    current_session: Option<Session>,
}

impl Client {
    /// 创建新的客户端实例
    fn new(config: AppConfig) -> Self {
        Self {
            config,
            current_session: None,
        }
    }

    /// 初始化客户端
    async fn init(&mut self) -> Result<()> {
        tracing::info!("{}", t!("client.connecting"));

        // 这里应该是连接服务器的逻辑
        // 暂时模拟连接成功
        tracing::info!("{}", t!("client.connected"));

        Ok(())
    }

    /// 连接到服务器
    async fn connect(&mut self, host: String, port: u16) -> Result<()> {
        self.config.server_url = host;
        self.config.server_port = port;
        
        tracing::info!(
            "{} {}:{}", 
            t!("client.connecting"), 
            self.config.server_url, 
            self.config.server_port
        );
        
        // 实际连接逻辑
        Ok(())
    }

    /// 列出会话
    async fn list_sessions(&self) -> Result<Vec<Session>> {
        tracing::info!("{}", t!("session.no_sessions"));
        
        // 返回模拟数据
        Ok(vec![])
    }

    /// 列出监听器
    async fn list_listeners(&self) -> Result<Vec<ListenerConfig>> {
        tracing::info!("{}", t!("listener.no_listeners"));
        
        // 返回模拟数据
        Ok(vec![])
    }

    /// 生成载荷
    async fn generate_payload(&self, config: PayloadConfig) -> Result<()> {
        tracing::info!("{}", t!("payload.generating"));
        
        // 实际生成逻辑
        tracing::info!("{}", t!("payload.generated"));
        
        Ok(())
    }

    /// 显示欢迎信息
    fn print_welcome(&self) {
        execute!(stdout(), Clear(ClearType::All)).unwrap();
        println!("{}", t!("common.welcome"));
        println!("{}", t!("common.version", version = "0.1.0"));
        println!();
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    // 解析命令行参数
    let args = Args::parse();

    // 设置语言
    let locale = match args.language.as_str() {
        "zh-CN" => Locale::zh_CN,
        "en" => Locale::en,
        _ => Locale::zh_CN,
    };
    rust_i18n::set_locale(&locale);

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

    // 创建客户端配置
    let config = AppConfig {
        server_url: args.server.clone(),
        server_port: args.port,
        language: args.language.clone(),
        log_level: args.log_level.clone(),
    };

    // 创建客户端实例
    let mut client = Client::new(config);
    
    // 显示欢迎信息
    client.print_welcome();

    // 初始化客户端
    client.init().await?;

    // 处理子命令
    match args.command {
        Some(Commands::Connect { host, port }) => {
            let host = host.unwrap_or(args.server);
            let port = port.unwrap_or(args.port);
            client.connect(host, port).await?;
        }
        
        Some(Commands::Session { session_command }) => {
            match session_command {
                SessionCommands::List => {
                    let sessions = client.list_sessions().await?;
                    println!("{}", t!("session.list_header"));
                    for session in sessions {
                        println!("  {} - {}@{}", session.name, session.username, session.hostname);
                    }
                }
                
                SessionCommands::Use { session_id } => {
                    println!("{}: {}", t!("session.select"), session_id);
                }
                
                SessionCommands::Kill { session_id } => {
                    println!("{}: {}", t!("session.kill_confirm"), session_id);
                }
                
                SessionCommands::Info { session_id } => {
                    println!("{}: {}", t!("session.session_info"), session_id);
                }
            }
        }
        
        Some(Commands::Listener { listener_command }) => {
            match listener_command {
                ListenerCommands::List => {
                    let listeners = client.list_listeners().await?;
                    println!("{}", t!("listener.list_header"));
                }
                
                ListenerCommands::Create { listener_type, host, port } => {
                    println!("{} {} {}:{} ({})", 
                        t!("listener.create"), 
                        t!("listener.types.mtls"), 
                        host, 
                        port, 
                        listener_type
                    );
                }
                
                ListenerCommands::Remove { name } => {
                    println!("{}: {}", t!("listener.remove_confirm"), name);
                }
            }
        }
        
        Some(Commands::Generate { generate_command }) => {
            match generate_command {
                GenerateCommands::Payload { payload_type, os, arch, c2_transport } => {
                    let config = PayloadConfig {
                        name: format!("{}_{}_{}", os, arch, payload_type),
                        payload_type,
                        os,
                        arch,
                        format: "exe".to_string(),
                        c2_transport,
                    };
                    client.generate_payload(config).await?;
                }
                
                GenerateCommands::List => {
                    println!("{}", t!("payload.list_header"));
                    println!("{}", t!("payload.no_payloads"));
                }
            }
        }
        
        Some(Commands::Beacon { beacon_command }) => {
            match beacon_command {
                BeaconCommands::List => {
                    println!("{}", t!("beacon.list_header"));
                    println!("{}", t!("beacon.no_beacons"));
                }
                
                BeaconCommands::Interval { seconds } => {
                    println!("{} {} {}", t!("beacon.set_interval"), t!("command.name"), seconds);
                }
                
                BeaconCommands::Jitter { percent } => {
                    println!("{} {} {}%", t!("beacon.set_jitter"), t!("command.name"), percent);
                }
            }
        }
        
        None => {
            // 没有子命令时，进入交互模式
            println!("{}", t!("help.usage"));
            println!("  --help    {}", t!("help.description"));
        }
    }

    // 禁用原始模式
    disable_raw_mode()?;

    Ok(())
}
