# Origami LSP Server Makefile

# 变量定义
BINARY_NAME=origami-lsp
BUILD_DIR=build
INSTALL_PATH=/usr/local/bin
VSCODE_EXT_DIR=vscode-extension
VSIX_FILE=$(VSCODE_EXT_DIR)/origami-language-support-1.0.0.vsix

# 默认目标
.PHONY: all
all: build

# 构建目标
.PHONY: build
build:
	@echo "构建 Origami LSP 服务器..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 安装目标
.PHONY: install
install: build
	@echo "安装 LSP 服务器到 $(INSTALL_PATH)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@echo "安装完成"

# 卸载目标
.PHONY: uninstall
uninstall:
	@echo "卸载 LSP 服务器..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "卸载完成"

# 清理目标
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 测试目标
.PHONY: test
test:
	@echo "运行测试..."
	go test ./...

# 运行目标
.PHONY: run
run: build
	@echo "启动 LSP 服务器..."
	$(BUILD_DIR)/$(BINARY_NAME)

# 开发模式运行
.PHONY: dev
dev:
	@echo "开发模式启动 LSP 服务器..."
	go run main.go

# 检查代码格式
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "检查代码..."
	go vet ./...

# 完整检查
.PHONY: check
check: fmt vet test
	@echo "所有检查完成"

# VS Code 扩展相关目标
.PHONY: vscode-build
vscode-build:
	@echo "构建 VS Code 扩展..."
	cd $(VSCODE_EXT_DIR) && npm install && npm run compile
	@echo "VS Code 扩展构建完成"

.PHONY: vscode-package
vscode-package: vscode-build
	@echo "打包 VS Code 扩展..."
	cd $(VSCODE_EXT_DIR) && npm run package
	@echo "VS Code 扩展打包完成: $(VSIX_FILE)"

.PHONY: vscode-install
vscode-install: vscode-package
	@echo "安装 VS Code 扩展..."
	code --install-extension $(VSIX_FILE)
	@echo "VS Code 扩展安装完成"

.PHONY: vscode-uninstall
vscode-uninstall:
	@echo "卸载 VS Code 扩展..."
	code --uninstall-extension origami-lang.origami-language-support
	@echo "VS Code 扩展卸载完成"

.PHONY: vscode-clean
vscode-clean:
	@echo "清理 VS Code 扩展构建文件..."
	cd $(VSCODE_EXT_DIR) && rm -rf out node_modules *.vsix
	@echo "VS Code 扩展清理完成"

# 完整安装（LSP 服务器 + VS Code 扩展）
.PHONY: install-all
install-all: install vscode-install
	@echo "完整安装完成（LSP 服务器 + VS Code 扩展）"

# 完整卸载
.PHONY: uninstall-all
uninstall-all: uninstall vscode-uninstall
	@echo "完整卸载完成"

# 完整清理
.PHONY: clean-all
clean-all: clean vscode-clean
	@echo "完整清理完成"

# 帮助信息
.PHONY: help
help:
	@echo "Origami LSP Server Makefile"
	@echo ""
	@echo "LSP 服务器目标:"
	@echo "  build     - 构建 LSP 服务器"
	@echo "  install   - 安装 LSP 服务器到系统路径"
	@echo "  uninstall - 卸载 LSP 服务器"
	@echo "  clean     - 清理 LSP 服务器构建文件"
	@echo "  test      - 运行测试"
	@echo "  run       - 构建并运行 LSP 服务器"
	@echo "  dev       - 开发模式运行 LSP 服务器"
	@echo "  fmt       - 格式化代码"
	@echo "  vet       - 代码检查"
	@echo "  check     - 完整检查 (fmt + vet + test)"
	@echo ""
	@echo "VS Code 扩展目标:"
	@echo "  vscode-build     - 构建 VS Code 扩展"
	@echo "  vscode-package   - 打包 VS Code 扩展"
	@echo "  vscode-install   - 安装 VS Code 扩展"
	@echo "  vscode-uninstall - 卸载 VS Code 扩展"
	@echo "  vscode-clean     - 清理 VS Code 扩展构建文件"
	@echo ""
	@echo "组合目标:"
	@echo "  install-all   - 安装 LSP 服务器和 VS Code 扩展"
	@echo "  uninstall-all - 卸载 LSP 服务器和 VS Code 扩展"
	@echo "  clean-all     - 清理所有构建文件"
	@echo "  help          - 显示此帮助信息"