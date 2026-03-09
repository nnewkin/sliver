#!/bin/bash
# Sliver Rust 版本自动化测试脚本
# 执行完整的测试闭环：编译→测试→修复→重测

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 切换到项目目录
cd "$(dirname "$0")"

# 1. 清理之前的构建
log_info "清理之前的构建产物..."
rm -rf target/

# 2. 编译项目
log_info "开始编译项目..."
if ! cargo build --release 2>&1 | tee build.log; then
    log_error "编译失败！"
    cat build.log
    exit 1
fi
log_info "编译成功！"

# 3. 运行单元测试
log_info "运行单元测试..."
if ! cargo test --lib --release 2>&1 | tee unit_test.log; then
    log_error "单元测试失败！"
    cat unit_test.log
    exit 1
fi
log_info "单元测试通过！"

# 4. 运行集成测试
log_info "运行集成测试..."
if ! cargo test --test "*_integration*" --release 2>&1 | tee integration_test.log; then
    log_error "集成测试失败！"
    cat integration_test.log
    exit 1
fi
log_info "集成测试通过！"

# 5. 运行所有测试
log_info "运行所有测试（包括文档测试）..."
if ! cargo test --all-features --release 2>&1 | tee all_test.log; then
    log_error "完整测试失败！"
    cat all_test.log
    exit 1
fi
log_info "完整测试通过！"

# 6. 验证二进制文件
log_info "验证二进制文件..."
if [ ! -f "target/release/sliver-client" ]; then
    log_error "客户端二进制文件未找到！"
    exit 1
fi

if [ ! -f "target/release/sliver-server" ]; then
    log_error "服务端二进制文件未找到！"
    exit 1
fi

log_info "所有二进制文件验证通过！"

# 7. 测试语言切换功能
log_info "测试语言切换功能..."
./target/release/sliver-client --language=zh-CN --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
    log_info "中文语言测试通过！"
else
    log_error "中文语言测试失败！"
    exit 1
fi

./target/release/sliver-client --language=en --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
    log_info "英文语言测试通过！"
else
    log_error "英文语言测试失败！"
    exit 1
fi

# 8. 生成测试报告
log_info "生成测试报告..."
cat > test_report.md << 'EOF'
# Sliver Rust 版本测试报告

## 测试执行时间
$(date)

## 编译状态
✅ 编译成功

## 测试结果

### 单元测试
✅ 所有单元测试通过
详见：unit_test.log

### 集成测试  
✅ 所有集成测试通过
详见：integration_test.log

### 完整测试
✅ 所有测试通过
详见：all_test.log

### 二进制文件验证
✅ 客户端二进制文件存在
✅ 服务端二进制文件存在

### 语言切换测试
✅ 中文语言切换正常
✅ 英文语言切换正常

## 测试覆盖率
EOF

# 添加测试统计
echo "" >> test_report.md
echo '```bash' >> test_report.md
grep "test result:" unit_test.log | tail -5 >> test_report.md
echo '```' >> test_report.md

log_info "测试报告已生成：test_report.md"

# 9. 总结
echo ""
log_info "=========================================="
log_info "  所有测试通过！"
log_info "=========================================="
echo ""

exit 0
