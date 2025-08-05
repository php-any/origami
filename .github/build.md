# GitHub Actions 构建说明

本项目使用 GitHub Actions 自动构建和发布二进制文件。

## 工作流程

### 1. 构建测试 (build.yml)
- **触发条件**: 推送到主分支、Pull Request、手动触发
- **功能**: 编译测试，确保代码可以正常构建
- **平台**: Linux, Windows, macOS (仅 amd64)
- **产物**: 临时构建文件（保留7天）

### 2. 发布构建 (release.yml)
- **触发条件**: 推送标签（如 `v1.0.0`）、手动触发
- **功能**: 构建多平台二进制文件并发布到 GitHub Releases
- **平台**: 
  - Linux (amd64, arm64)
  - Windows (amd64)
  - macOS (amd64, arm64)

## 如何发布新版本

1. 确保代码已经测试完毕
2. 创建并推送标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions 会自动：
   - 构建所有平台的二进制文件
   - 创建 GitHub Release
   - 上传构建产物到 Release

## 构建产物说明

- `origami-linux-amd64`: Linux x64 版本
- `origami-linux-arm64`: Linux ARM64 版本
- `origami-darwin-amd64`: macOS Intel 版本
- `origami-darwin-arm64`: macOS Apple Silicon 版本
- `origami-windows-amd64.exe`: Windows x64 版本

## 手动触发构建

可以在 GitHub 仓库的 Actions 页面手动触发构建：
1. 进入 Actions 页面
2. 选择对应的工作流
3. 点击 "Run workflow" 按钮