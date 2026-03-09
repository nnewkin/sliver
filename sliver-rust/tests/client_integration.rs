// Sliver 客户端集成测试
// 测试 CLI 客户端的核心功能，包括会话、监听器、载荷生成等

use rstest::*;
use std::env;

// 测试客户端帮助命令
#[test]
fn test_client_help() {
    use std::process::Command;
    
    // 从环境变量获取项目根目录
    let project_root = env::var("CARGO_MANIFEST_DIR").unwrap_or_else(|_| ".".to_string());
    
    // 构建二进制文件路径
    let client_path = if cfg!(debug_assertions) {
        format!("{}/../../target/debug/sliver-client", project_root)
    } else {
        format!("{}/../../target/release/sliver-client", project_root)
    };
    
    // 检查文件是否存在，不存在则跳过测试
    if !std::path::Path::new(&client_path).exists() {
        println!("警告：客户端二进制文件不存在: {}，跳过测试", client_path);
        return;
    }
    
    let output = Command::new(&client_path)
        .arg("--help")
        .output()
        .expect("执行客户端帮助命令失败");
    
    assert!(output.status.success());
    let help_text = String::from_utf8_lossy(&output.stdout);
    
    // 验证中文帮助信息
    assert!(help_text.contains("sliver-client"));
    assert!(help_text.contains("会话管理"));
    assert!(help_text.contains("监听器管理"));
    assert!(help_text.contains("载荷生成"));
}

// 测试语言切换功能
#[test]
fn test_language_switching() {
    use std::process::Command;
    
    // 测试中文
    let output_zh = Command::new("./target/debug/sliver-client")
        .arg("--language=zh-CN")
        .arg("session")
        .arg("list")
        .output()
        .expect("执行中文会话列表命令失败");
    
    assert!(output_zh.status.success());
    let output_zh_text = String::from_utf8_lossy(&output_zh.stdout);
    assert!(output_zh_text.contains("会话列表"));
    
    // 测试英文
    let output_en = Command::new("./target/debug/sliver-client")
        .arg("--language=en")
        .arg("session")
        .arg("list")
        .output()
        .expect("执行英文会话列表命令失败");
    
    assert!(output_en.status.success());
    let output_en_text = String::from_utf8_lossy(&output_en.stdout);
    assert!(output_en_text.contains("Session List"));
}

// 测试会话命令
#[test]
fn test_session_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("session")
        .arg("list")
        .output()
        .expect("执行会话列表命令失败");
    
    assert!(output.status.success());
}

// 测试监听器命令
#[test]
fn test_listener_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("listener")
        .arg("list")
        .output()
        .expect("执行监听器列表命令失败");
    
    assert!(output.status.success());
}

// 测试文件操作命令
#[test]
fn test_file_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("file")
        .arg("list")
        .arg("--path")
        .arg(".")
        .output()
        .expect("执行文件列表命令失败");
    
    assert!(output.status.success());
}

// 测试进程命令
#[test]
fn test_process_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("process")
        .arg("list")
        .output()
        .expect("执行进程列表命令失败");
    
    assert!(output.status.success());
}

// 测试网络命令
#[test]
fn test_network_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("network")
        .arg("scan")
        .arg("--target")
        .arg("localhost")
        .arg("--ports")
        .arg("80,443")
        .output()
        .expect("执行网络扫描命令失败");
    
    assert!(output.status.success());
}

// 测试权限命令
#[test]
fn test_privilege_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("privilege")
        .arg("impersonate")
        .arg("--username")
        .arg("test")
        .output()
        .expect("执行权限命令失败");
    
    assert!(output.status.success());
}

// 测试载荷生成命令
#[test]
fn test_payload_commands() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-client")
        .arg("generate")
        .arg("list")
        .output()
        .expect("执行载荷列表命令失败");
    
    assert!(output.status.success());
}
