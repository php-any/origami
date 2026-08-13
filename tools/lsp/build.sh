#!/bin/bash

# Origami LSP 构建脚本
# 默认构建到上层目录的 bin/zy-lsp

set -e  # 遇到错误立即退出

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BIN_DIR="$PROJECT_ROOT/bin"
OUTPUT_FILE="$BIN_DIR/zy-lsp"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Origami LSP 构建脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -o, --output   指定输出文件路径 (默认: $OUTPUT_FILE)"
    echo "  -c, --clean    构建前清理输出目录"
    echo "  -v, --verbose  显示详细构建信息"
    echo "  --version      显示版本信息"
    echo ""
    echo "示例:"
    echo "  $0                    # 构建到默认位置"
    echo "  $0 -o /tmp/lsp        # 构建到指定位置"
    echo "  $0 -c -v              # 清理后详细构建"
}

# 显示版本信息
show_version() {
    echo "Origami LSP 构建脚本 v1.0.0"
}

# 默认参数
CLEAN=false
VERBOSE=false
VERSION=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        --version)
            show_version
            exit 0
            ;;
        *)
            print_error "未知参数: $1"
            echo "使用 $0 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

# 显示构建信息
print_info "开始构建 Origami LSP..."
print_info "项目根目录: $PROJECT_ROOT"
print_info "输出文件: $OUTPUT_FILE"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    print_error "未找到 Go 环境，请先安装 Go"
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3)
print_info "Go 版本: $GO_VERSION"

# 清理输出目录
if [ "$CLEAN" = true ]; then
    print_info "清理输出目录..."
    if [ -d "$BIN_DIR" ]; then
        rm -rf "$BIN_DIR"
        print_success "已清理输出目录: $BIN_DIR"
    else
        print_warning "输出目录不存在，无需清理"
    fi
fi

# 创建输出目录
if [ ! -d "$BIN_DIR" ]; then
    mkdir -p "$BIN_DIR"
    print_info "已创建输出目录: $BIN_DIR"
fi

# 进入 LSP 目录
cd "$SCRIPT_DIR"

# 检查 go.mod 文件
if [ ! -f "go.mod" ]; then
    print_error "未找到 go.mod 文件，请确保在正确的目录中运行此脚本"
    exit 1
fi

# 下载依赖
print_info "下载 Go 模块依赖..."
if [ "$VERBOSE" = true ]; then
    go mod download
else
    go mod download > /dev/null 2>&1
fi
print_success "依赖下载完成"

# 构建应用
print_info "开始编译..."
BUILD_FLAGS="-ldflags=-s"  # 去除调试信息和符号表，减小文件大小

if [ "$VERBOSE" = true ]; then
    go build $BUILD_FLAGS -o "$OUTPUT_FILE" .
else
    go build $BUILD_FLAGS -o "$OUTPUT_FILE" . > /dev/null 2>&1
fi

# 检查构建结果
if [ $? -eq 0 ]; then
    print_success "构建完成!"
    
    # 显示文件信息
    if [ -f "$OUTPUT_FILE" ]; then
        FILE_SIZE=$(du -h "$OUTPUT_FILE" | cut -f1)
        print_info "输出文件: $OUTPUT_FILE"
        print_info "文件大小: $FILE_SIZE"
        
        # 检查文件是否可执行
        if [ -x "$OUTPUT_FILE" ]; then
            print_success "文件已设置为可执行"
        else
            print_warning "文件未设置为可执行，正在修复..."
            chmod +x "$OUTPUT_FILE"
            print_success "已设置文件为可执行"
        fi
        
        # 显示版本信息
        print_info "测试构建的二进制文件..."
        if "$OUTPUT_FILE" --version > /dev/null 2>&1; then
            VERSION_OUTPUT=$("$OUTPUT_FILE" --version 2>&1)
            print_success "版本信息: $VERSION_OUTPUT"
        else
            print_warning "无法获取版本信息"
        fi
        
    else
        print_error "构建失败：未找到输出文件"
        exit 1
    fi
else
    print_error "构建失败"
    exit 1
fi

print_success "Origami LSP 构建完成！"
print_info "可以使用以下命令运行:"
print_info "  $OUTPUT_FILE [选项]"
print_info "  $OUTPUT_FILE --help  # 查看帮助信息"
