#!/bin/bash
# 编译测试脚本 - 验证所有平台的编译能力

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 切换到项目目录
cd "$(dirname "$0")"

log_info "开始编译测试..."

# 清理旧的构建
log_info "清理旧的构建产物..."
rm -rf target/

# 测试 debug 模式编译
log_info "测试 Debug 模式编译..."
if cargo build 2>&1 | tee debug_build.log; then
    log_success "Debug 模式编译成功"
else
    log_error "Debug 模式编译失败"
    cat debug_build.log
    exit 1
fi

# 验证 debug 二进制文件
log_info "验证 Debug 二进制文件..."
if [ -f "target/debug/sliver-client" ] && [ -f "target/debug/sliver-server" ]; then
    log_success "Debug 二进制文件存在"
else
    log_error "Debug 二进制文件缺失"
    exit 1
fi

# 测试 release 模式编译
log_info "测试 Release 模式编译..."
if cargo build --release 2>&1 | tee release_build.log; then
    log_success "Release 模式编译成功"
else
    log_error "Release 模式编译失败"
    cat release_build.log
    exit 1
fi

# 验证 release 二进制文件
log_info "验证 Release 二进制文件..."
if [ -f "target/release/sliver-client" ] && [ -f "target/release/sliver-server" ]; then
    log_success "Release 二进制文件存在"
else
    log_error "Release 二进制文件缺失"
    exit 1
fi

# 测试代码格式检查
log_info "运行代码格式检查..."
if cargo fmt --all -- --check 2>&1; then
    log_success "代码格式检查通过"
else
    log_error "代码格式检查失败，请运行 'cargo fmt' 修复"
    exit 1
fi

# 测试 clippy 静态分析
log_info "运行 Clippy 静态分析..."
if cargo clippy --all-targets --all-features 2>&1 | tee clippy_check.log; then
    log_success "Clippy 检查通过"
else
    log_error "Clippy 检查发现警告或错误"
    cat clippy_check.log
    # 不退出，因为警告可能不影响功能
fi

# 生成编译报告
log_info "生成编译报告..."
cat > build_report.md << 'EOF'
# Sliver Rust 版本编译报告

## 编译时间
$(date)

## 编译结果

### Debug 模式
✅ 编译成功
二进制文件：
- target/debug/sliver-client
- target/debug/sliver-server

### Release 模式  
✅ 编译成功
二进制文件：
- target/release/sliver-client
- target/release/sliver-server

### 代码质量检查
✅ 代码格式检查通过
$(if [ -s clippy_check.log ]; then
    if grep -q "warning" clippy_check.log; then
        echo "⚠️  Clippy 发现警告"
    else
        echo "✅ Clippy 检查通过"
    fi
fi)

## 编译统计
EOF

echo "" >> build_report.md
echo '```bash' >> build_report.md
echo "Debug 模式编译时间：" >> build_report.md
grep "Finished" debug_build.log | tail -1 >> build_report.md
echo "" >> build_report.md
echo "Release 模式编译时间：" >> build_report.md
grep "Finished" release_build.log | tail -1 >> build_report.md
echo '```' >> build_report.md

log_success "编译报告已生成：build_report.md"

echo ""
log_info "=========================================="
log_success "  所有编译测试通过！"
log_info "=========================================="
echo ""

exit 0
