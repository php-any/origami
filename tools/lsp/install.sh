#!/bin/bash

# Origami LSP 一键安装脚本
# 安装 LSP 服务器和 VS Code 扩展

set -e

echo "🚀 Origami LSP 一键安装脚本"
echo "================================"

# 检查依赖
check_dependencies() {
    echo "📋 检查依赖..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo "❌ 错误: 未找到 Go，请先安装 Go"
        exit 1
    fi
    echo "✅ Go: $(go version)"
    
    # 检查 VS Code
    if ! command -v code &> /dev/null; then
        echo "⚠️  警告: 未找到 VS Code，将跳过扩展安装"
        INSTALL_VSCODE=false
    else
        echo "✅ VS Code: $(code --version | head -n1)"
        INSTALL_VSCODE=true
    fi
    
    # 检查 Node.js (VS Code 扩展需要)
    if [ "$INSTALL_VSCODE" = true ]; then
        if ! command -v npm &> /dev/null; then
            echo "⚠️  警告: 未找到 npm，将跳过 VS Code 扩展安装"
            INSTALL_VSCODE=false
        else
            echo "✅ npm: $(npm --version)"
        fi
    fi
}

# 安装 LSP 服务器
install_lsp_server() {
    echo ""
    echo "🔧 安装 LSP 服务器..."
    make install
    
    # 验证安装
    if command -v origami-lsp &> /dev/null; then
        echo "✅ LSP 服务器安装成功: $(which origami-lsp)"
    else
        echo "❌ LSP 服务器安装失败"
        exit 1
    fi
}

# 安装 VS Code 扩展
install_vscode_extension() {
    if [ "$INSTALL_VSCODE" = true ]; then
        echo ""
        echo "🎨 安装 VS Code 扩展..."
        make vscode-install
        
        # 验证安装
        if code --list-extensions | grep -q "origami-lang.origami-language-support"; then
            echo "✅ VS Code 扩展安装成功"
        else
            echo "❌ VS Code 扩展安装失败"
            exit 1
        fi
    else
        echo ""
        echo "⏭️  跳过 VS Code 扩展安装"
    fi
}

# 创建工作区配置
create_workspace_config() {
    echo ""
    echo "⚙️  创建工作区配置..."
    
    # 检查是否已存在 .vscode 目录
    if [ ! -d "../../.vscode" ]; then
        mkdir -p "../../.vscode"
    fi
    
    # 创建 settings.json（如果不存在）
    if [ ! -f "../../.vscode/settings.json" ]; then
        cat > "../../.vscode/settings.json" << 'EOF'
{
    "origami.lsp.enabled": true,
    "origami.lsp.serverPath": "/usr/local/bin/origami-lsp",
    "origami.lsp.trace": "verbose",
    "files.associations": {
        "*.cjp": "origami",
        "*.origami": "origami"
    },
    "editor.quickSuggestions": {
        "other": true,
        "comments": false,
        "strings": true
    },
    "editor.suggest.insertMode": "replace",
    "editor.acceptSuggestionOnCommitCharacter": false
}
EOF
        echo "✅ 工作区配置已创建: ../../.vscode/settings.json"
    else
        echo "ℹ️  工作区配置已存在，跳过创建"
    fi
}

# 显示安装结果
show_results() {
    echo ""
    echo "🎉 安装完成！"
    echo "=============="
    
    echo ""
    echo "📦 已安装组件:"
    echo "  • LSP 服务器: $(which origami-lsp)"
    
    if [ "$INSTALL_VSCODE" = true ]; then
        echo "  • VS Code 扩展: origami-lang.origami-language-support"
    fi
    
    echo ""
    echo "🚀 快速开始:"
    echo "  1. 在 VS Code 中打开 .cjp 或 .origami 文件"
    echo "  2. 享受语法高亮、代码补全和悬停信息"
    echo "  3. 使用 Ctrl+Shift+P 打开命令面板，搜索 'Origami'"
    
    echo ""
    echo "📚 文档:"
    echo "  • 安装指南: INSTALLATION.md"
    echo "  • VS Code 指南: VSCODE_INSTALLATION.md"
    
    echo ""
    echo "🔧 管理命令:"
    echo "  • 卸载: make uninstall-all"
    echo "  • 重新安装: make install-all"
    echo "  • 查看帮助: make help"
}

# 主函数
main() {
    check_dependencies
    install_lsp_server
    install_vscode_extension
    create_workspace_config
    show_results
}

# 运行主函数
main "$@"