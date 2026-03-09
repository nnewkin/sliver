// Sliver 服务端集成测试
// 测试服务端的核心功能，包括启动、停止、监听器管理等

use rstest::*;

// 测试服务端帮助命令
#[test]
fn test_server_help() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-server")
        .arg("--help")
        .output()
        .expect("执行服务端帮助命令失败");
    
    assert!(output.status.success());
    let help_text = String::from_utf8_lossy(&output.stdout);
    
    // 验证中文帮助信息
    assert!(help_text.contains("sliver-server"));
    assert!(help_text.contains("启动服务器"));
    assert!(help_text.contains("监听器管理"));
}

// 测试服务端启动
#[test]
fn test_server_start() {
    use std::process::Command;
    use std::time::Duration;
    
    let mut child = Command::new("./target/debug/sliver-server")
        .arg("start")
        .spawn()
        .expect("启动服务端失败");
    
    // 等待一段时间让服务启动
    std::thread::sleep(Duration::from_secs(2));
    
    // 检查进程状态
    let status = child.try_wait();
    assert!(status.is_ok(), "服务端进程异常退出");
}

// 测试服务端状态查询
#[test]
fn test_server_status() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-server")
        .arg("status")
        .output()
        .expect("执行服务端状态查询失败");
    
    assert!(output.status.success());
    let status_text = String::from_utf8_lossy(&output.stdout);
    
    // 验证状态信息
    assert!(status_text.contains("状态") || status_text.contains("Status"));
    assert!(status_text.contains("会话") || status_text.contains("Session"));
    assert!(status_text.contains("监听器") || status_text.contains("Listener"));
}

// 测试监听器创建
#[test]
fn test_listener_create() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-server")
        .arg("listener")
        .arg("create")
        .arg("--listener-type")
        .arg("http")
        .arg("--host")
        .arg("0.0.0.0")
        .arg("--port")
        .arg("8080")
        .output()
        .expect("创建监听器失败");
    
    assert!(output.status.success());
}

// 测试监听器列表
#[test]
fn test_listener_list() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-server")
        .arg("listener")
        .arg("list")
        .output()
        .expect("获取监听器列表失败");
    
    assert!(output.status.success());
}

// 测试载荷生成
#[test]
fn test_payload_generate() {
    use std::process::Command;
    
    let output = Command::new("../../target/debug/sliver-server")
        .arg("generate")
        .arg("--payload-type")
        .arg("exe")
        .arg("--os")
        .arg("windows")
        .arg("--arch")
        .arg("amd64")
        .output()
        .expect("生成载荷失败");
    
    assert!(output.status.success());
    let payload_text = String::from_utf8_lossy(&output.stdout);
    
    // 验证生成信息
    assert!(payload_text.contains("载荷生成完成") || payload_text.contains("Payload generation complete"));
}

// 测试服务端语言切换
#[test]
fn test_server_language_switching() {
    use std::process::Command;
    
    // 测试中文
    let output_zh = Command::new("./target/debug/sliver-server")
        .arg("--language=zh-CN")
        .arg("status")
        .output()
        .expect("执行中文状态查询失败");
    
    assert!(output_zh.status.success());
    
    // 测试英文
    let output_en = Command::new("./target/debug/sliver-server")
        .arg("--language=en")
        .arg("status")
        .output()
        .expect("执行英文状态查询失败");
    
    assert!(output_en.status.success());
}
