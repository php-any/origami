# 安装指南

本指南将帮助您安装和配置折言(Origami)语言环境。

## 系统要求

### 必需条件

- **Go 1.18+** - 用于编译折言语言
- **Git** - 用于克隆代码仓库
- **操作系统**: Linux, macOS, Windows

### 推荐配置

- **内存**: 4GB+
- **磁盘空间**: 1GB+
- **网络**: 稳定的网络连接（用于下载依赖）

## 安装步骤

### 1. 安装 Go

#### Linux/macOS

```bash
# 下载并安装 Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### macOS (使用 Homebrew)

```bash
brew install go
```

#### Windows

1. 访问 [Go 官网](https://golang.org/dl/)
2. 下载 Windows 安装包
3. 运行安装程序并按照提示完成安装

### 2. 验证 Go 安装

```bash
go version
```

应该显示类似输出：

```
go version go1.21.0 linux/amd64
```

### 3. 克隆折言仓库

```bash
git clone https://github.com/php-any/origami.git
cd origami
```

### 4. 编译折言语言

```bash
go build -o origami origami.go
```

### 5. 验证安装

```bash
./origami
```

应该显示帮助信息：

```
折言(origami-lang) - 融合型脚本语言

用法: ./origami <脚本路径>
...
```

## 运行测试

验证安装是否成功：

```bash
./origami tests/run_tests.cjp
```

如果看到大量测试通过信息，说明安装成功。

## 配置开发环境

### 1. 设置 PATH（可选）

将折言可执行文件添加到系统 PATH：

#### Linux/macOS

```bash
# 创建符号链接
sudo ln -s $(pwd)/origami /usr/local/bin/origami

# 验证
origami
```

#### Windows

1. 将 `origami.exe` 复制到系统 PATH 目录
2. 或在当前目录使用 `.\origami`

### 2. IDE 配置

#### VS Code

1. 安装 PHP 扩展
2. 配置文件关联：
   ```json
   {
     "files.associations": {
       "*.cjp": "php"
     }
   }
   ```

#### IntelliJ IDEA

1. 安装 PHP 插件
2. 配置 `.cjp` 文件类型为 PHP

## 故障排除

### 常见问题

#### 1. Go 版本过低

**错误**: `go: requires go1.18 or later`
**解决方案**: 升级 Go 到 1.18 或更高版本

#### 2. 依赖下载失败

**错误**: `go: module lookup disabled`
**解决方案**:

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

#### 3. 权限问题

**错误**: `permission denied`
**解决方案**:

```bash
chmod +x origami
```

#### 4. 编译错误

**错误**: 编译时出现各种错误
**解决方案**:

```bash
# 清理并重新编译
go clean
go mod tidy
go build -o origami origami.go
```

### 获取帮助

如果遇到其他问题：

1. 查看 [GitHub Issues](https://github.com/php-any/origami/issues)
2. 加入 [讨论群](https://github.com/php-any/origami#-讨论群)
3. 提交新的 Issue 描述问题

## 下一步

安装完成后，建议：

1. 阅读 [快速开始](quickstart.md) 创建第一个程序
2. 学习 [基础语法](syntax.md) 了解语言特性
3. 查看 [示例代码](https://github.com/php-any/origami/tree/main/tests) 学习最佳实践
