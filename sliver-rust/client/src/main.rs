// Sliver 客户端 - CLI 命令行界面
// 包含所有用户交互逻辑，使用 rust-i18n 实现多语言支持

use anyhow::Result;
use clap::{Parser, Subcommand};
use crossterm::{
    execute,
    terminal::{disable_raw_mode, Clear, ClearType},
};
use rust_i18n::t;
use sliver_shared::{Session, ListenerConfig, PayloadConfig, AppConfig};
use std::io::stdout;

pub mod commands;
use commands::filesystem;
use commands::process;
use commands::network;
use commands::privilege;
use commands::interactive;

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

    /// 文件操作
    File {
        /// 子命令
        #[command(subcommand)]
        file_command: FileCommands,
    },

    /// 进程管理
    Process {
        /// 子命令
        #[command(subcommand)]
        process_command: ProcessCommands,
    },

    /// 网络操作
    Network {
        /// 子命令
        #[command(subcommand)]
        network_command: NetworkCommands,
    },

    /// 权限管理
    Privilege {
        /// 子命令
        #[command(subcommand)]
        privilege_command: PrivilegeCommands,
    },

    /// 执行命令
    Execute {
        /// 要执行的命令
        #[arg(short, long)]
        command: String,
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

/// 文件操作命令
#[derive(Subcommand, Debug)]
enum FileCommands {
    /// 上传文件
    Upload {
        /// 本地文件路径
        #[arg(short, long)]
        local: String,

        /// 远程文件路径
        #[arg(short, long)]
        remote: String,
    },

    /// 下载文件
    Download {
        /// 远程文件路径
        #[arg(short, long)]
        remote: String,

        /// 本地文件路径
        #[arg(short, long)]
        local: String,
    },

    /// 列出文件
    List {
        /// 目录路径
        #[arg(short, long)]
        path: Option<String>,
    },

    /// 删除文件
    Delete {
        /// 文件路径
        #[arg(short, long)]
        path: String,

        /// 确认删除
        #[arg(short, long)]
        confirm: bool,
    },
}

/// 进程管理命令
#[derive(Subcommand, Debug)]
enum ProcessCommands {
    /// 列出进程
    List,

    /// 终止进程
    Kill {
        /// 进程ID
        #[arg(short, long)]
        pid: u32,

        /// 确认终止
        #[arg(short, long)]
        confirm: bool,
    },

    /// 进程迁移
    Migrate {
        /// 目标进程ID
        #[arg(short, long)]
        target_pid: u32,
    },

    /// 进程注入
    Inject {
        /// 载荷路径
        #[arg(short, long)]
        payload: String,
    },
}

/// 网络操作命令
#[derive(Subcommand, Debug)]
enum NetworkCommands {
    /// 端口扫描
    Scan {
        /// 目标主机
        #[arg(short, long)]
        target: String,

        /// 端口列表
        #[arg(short, long)]
        ports: String,
    },

    /// 设置 SOCKS 代理
    Socks {
        /// 本地端口
        #[arg(short, long)]
        local_port: u16,

        /// 远程主机
        #[arg(short, long)]
        remote_host: String,

        /// 远程端口
        #[arg(short, long)]
        remote_port: u16,
    },

    /// 反向端口转发
    Rportfwd {
        /// 本地端口
        #[arg(short, long)]
        local_port: u16,

        /// 远程端口
        #[arg(short, long)]
        remote_port: u16,
    },

    /// TCP 枢纽
    Tcppivot {
        /// 本地端口
        #[arg(short, long)]
        local_port: u16,

        /// 远程主机
        #[arg(short, long)]
        remote_host: String,

        /// 远程端口
        #[arg(short, long)]
        remote_port: u16,
    },
}

/// 权限管理命令
#[derive(Subcommand, Debug)]
enum PrivilegeCommands {
    /// 模拟用户
    Impersonate {
        /// 用户名
        #[arg(short, long)]
        username: String,
    },

    /// 获取 SYSTEM 权限
    Getsystem,

    /// 创建令牌
    MakeToken {
        /// 用户名
        #[arg(short, long)]
        username: String,

        /// 密码
        #[arg(short, long)]
        password: String,
    },

    /// 恢复原始令牌
    Rev2self,
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
    async fn generate_payload(&self, _config: PayloadConfig) -> Result<()> {
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
                    let _listeners = client.list_listeners().await?;
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

        Some(Commands::File { file_command }) => {
            match file_command {
                FileCommands::Upload { local, remote } => {
                    let result = filesystem::upload_file(&local, &remote).await?;
                    println!("{}", result.stdout);
                }

                FileCommands::Download { remote, local } => {
                    let result = filesystem::download_file(&remote, &local).await?;
                    println!("{}", result.stdout);
                }

                FileCommands::List { path } => {
                    let path = path.unwrap_or(".".to_string());
                    let result = filesystem::list_files(&path).await?;
                    println!("{}", result.stdout);
                }

                FileCommands::Delete { path, confirm } => {
                    let result = filesystem::delete_file(&path, confirm).await?;
                    println!("{}", result.stdout);
                }
            }
        }

        Some(Commands::Process { process_command }) => {
            match process_command {
                ProcessCommands::List => {
                    let result = process::list_processes().await?;
                    println!("{}", result.stdout);
                }

                ProcessCommands::Kill { pid, confirm } => {
                    let result = process::kill_process(pid, confirm).await?;
                    println!("{}", result.stdout);
                }

                ProcessCommands::Migrate { target_pid } => {
                    let result = process::migrate_process(0, target_pid).await?;
                    println!("{}", result.stdout);
                }

                ProcessCommands::Inject { payload } => {
                    let result = process::inject_process(0, &payload).await?;
                    println!("{}", result.stdout);
                }
            }
        }

        Some(Commands::Network { network_command }) => {
            match network_command {
                NetworkCommands::Scan { target, ports } => {
                    let port_list: Vec<u16> = ports.split(',')
                        .filter_map(|p| p.trim().parse().ok())
                        .collect();
                    let result = network::port_scan(&target, port_list).await?;
                    println!("{}", result.stdout);
                }

                NetworkCommands::Socks { local_port, remote_host, remote_port } => {
                    let result = network::setup_socks_proxy(local_port, &remote_host, remote_port).await?;
                    println!("{}", result.stdout);
                }

                NetworkCommands::Rportfwd { local_port, remote_port } => {
                    let result = network::setup_reverse_port_forward(local_port, remote_port).await?;
                    println!("{}", result.stdout);
                }

                NetworkCommands::Tcppivot { local_port, remote_host, remote_port } => {
                    let result = network::setup_tcp_pivot(local_port, &remote_host, remote_port).await?;
                    println!("{}", result.stdout);
                }
            }
        }

        Some(Commands::Privilege { privilege_command }) => {
            match privilege_command {
                PrivilegeCommands::Impersonate { username } => {
                    let result = privilege::impersonate_user(&username).await?;
                    println!("{}", result.stdout);
                }

                PrivilegeCommands::Getsystem => {
                    let result = privilege::get_system().await?;
                    println!("{}", result.stdout);
                }

                PrivilegeCommands::MakeToken { username, password } => {
                    let result = privilege::make_token(&username, &password).await?;
                    println!("{}", result.stdout);
                }

                PrivilegeCommands::Rev2self => {
                    let result = privilege::rev2self().await?;
                    println!("{}", result.stdout);
                }
            }
        }

        Some(Commands::Execute { command }) => {
            let result = interactive::execute_command(&command).await?;
            println!("{}", result.stdout);
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
