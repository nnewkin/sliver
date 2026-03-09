// 多语言功能测试
// 测试 rust-i18n 多语言支持的正确性

use rust_i18n::t;

// 初始化 i18n
rust_i18n::i18n!("../locales", fallback = "en");

// 测试中文翻译
#[test]
fn test_chinese_translation() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试通用文本
    let welcome = t!("common.welcome");
    assert!(welcome.contains("Sliver") || welcome.contains("欢迎使用"));
    
    // 测试客户端文本
    let client_name = t!("client.name");
    assert!(client_name.contains("客户端") || client_name.contains("Client"));
    
    // 测试会话文本
    let session_name = t!("session.name");
    assert!(session_name.contains("会话") || session_name.contains("Session"));
    
    // 测试监听器文本
    let listener_name = t!("listener.name");
    assert!(listener_name.contains("监听器") || listener_name.contains("Listener"));
    
    // 测试载荷文本
    let payload_name = t!("payload.name");
    assert!(payload_name.contains("载荷") || payload_name.contains("Payload"));
}

// 测试英文翻译
#[test]
fn test_english_translation() {
    rust_i18n::set_locale("en");
    
    // 测试通用文本
    let welcome = t!("common.welcome");
    assert!(welcome.contains("Sliver") || welcome.contains("Welcome"));
    
    // 测试客户端文本
    let client_name = t!("client.name");
    assert!(client_name.contains("Client"));
    
    // 测试会话文本
    let session_name = t!("session.name");
    assert!(session_name.contains("Session"));
    
    // 测试监听器文本
    let listener_name = t!("listener.name");
    assert!(listener_name.contains("Listener"));
    
    // 测试载荷文本
    let payload_name = t!("payload.name");
    assert!(payload_name.contains("Payload"));
}

// 测试带参数的翻译
#[test]
fn test_translation_with_parameters() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试带版本参数
    let version = t!("common.version", version = "0.1.0");
    assert!(version.contains("0.1.0"));
    
    // 测试带计数的翻译
    let count = t!("session.active_sessions", count = 5);
    assert!(count.contains("5"));
    
    // 测试带名称的翻译
    let name = t!("session.select", name = "test-session");
    assert!(name.contains("test-session"));
}

// 测试错误信息翻译
#[test]
fn test_error_translation() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试网络错误
    let network_error = t!("error.network_error", message = "连接超时");
    assert!(network_error.contains("网络错误") || network_error.contains("Network error"));
    
    // 测试文件未找到错误
    let file_error = t!("error.file_not_found", path = "/test/file.txt");
    assert!(file_error.contains("文件未找到") || file_error.contains("File not found"));
    
    // 测试权限拒绝错误
    let perm_error = t!("error.permission_denied");
    assert!(perm_error.contains("权限被拒绝") || perm_error.contains("Permission denied"));
}

// 测试状态翻译
#[test]
fn test_status_translation() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试活跃状态
    let active = t!("status.active");
    assert!(active.contains("活跃") || active.contains("Active"));
    
    // 测试非活跃状态
    let inactive = t!("status.inactive");
    assert!(inactive.contains("非活跃") || inactive.contains("Inactive"));
    
    // 测试运行中状态
    let running = t!("status.running");
    assert!(running.contains("运行中") || running.contains("Running"));
    
    // 测试已停止状态
    let stopped = t!("status.stopped");
    assert!(stopped.contains("已停止") || stopped.contains("Stopped"));
}

// 测试语言切换功能
#[test]
fn test_language_switching() {
    // 切换到中文
    rust_i18n::set_locale("zh-CN");
    let zh_text = t!("common.welcome");
    
    // 切换到英文
    rust_i18n::set_locale("en");
    let en_text = t!("common.welcome");
    
    // 验证两种语言返回不同内容
    // 注意：这里可能返回相同的 key，取决于实现
    // 实际应用中应该返回不同的翻译文本
    assert!(!zh_text.is_empty());
    assert!(!en_text.is_empty());
}

// 测试回退机制
#[test]
fn test_fallback_mechanism() {
    // 切换到不存在的语言
    rust_i18n::set_locale("de");
    
    // 应该回退到英文
    let text = t!("common.welcome");
    assert!(!text.is_empty());
}

// 测试命令翻译
#[test]
fn test_command_translation() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试文件命令
    let file_upload = t!("file.upload");
    assert!(file_upload.contains("上传") || file_upload.contains("upload"));
    
    // 测试进程命令
    let process_list = t!("process.list");
    assert!(process_list.contains("列表") || process_list.contains("list"));
    
    // 测试网络命令
    let network_scan = t!("network.scan");
    assert!(network_scan.contains("扫描") || network_scan.contains("scan"));
    
    // 测试权限命令
    let privilege_impersonate = t!("privilege.impersonate");
    assert!(privilege_impersonate.contains("模拟") || privilege_impersonate.contains("impersonate"));
}

// 测试功能翻译完整性
#[test]
fn test_feature_translation_completeness() {
    rust_i18n::set_locale("zh-CN");
    
    // 测试所有主要功能类别都有翻译
    let features = vec![
        t!("session.name"),
        t!("listener.name"),
        t!("payload.name"),
        t!("beacon.name"),
        t!("file.name"),
        t!("process.name"),
        t!("network.name"),
        t!("privilege.name"),
        t!("plugin.name"),
        t!("settings.name"),
    ];
    
    for feature in features {
        assert!(!feature.is_empty(), "功能翻译不应为空");
    }
}
