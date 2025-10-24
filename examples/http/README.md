# HTTP 扩展示例

本目录包含了 Origami 语言中 HTTP 功能的使用示例，展示了如何使用现有的 HTTP 服务器功能构建 Web 应用。

## 项目结构

```
examples/http/
├── go.mod                    # Go模块定义
├── main.go                   # 独立的Go程序入口
├── http.zy                   # HTTP服务器基础功能演示
└── README.md                 # 说明文档
```

## 功能特性

### HTTP 服务器基础功能

- 创建 HTTP 服务器
- 设置路由（GET、POST 等）
- 处理请求参数和路径参数
- JSON 响应处理
- 错误处理和状态码

### 中间件功能

- 全局日志记录
- CORS 处理
- 请求处理时间统计

## 使用方法

### 1. 进入 HTTP 例子目录

```bash
cd examples/http
```

### 2. 运行 HTTP 服务器示例

```bash
go run main.go
```

访问 `http://127.0.0.1:8080` 查看效果。

## API 端点

### HTTP 服务器 (端口 8080)

- `GET /` - 首页，显示服务器信息和可用端点
- `GET /user/{id}` - 获取用户信息
- `POST /user` - 创建用户（需要 name 和 email 参数）
- `GET /search` - 搜索功能（支持 keyword 和 page 参数）
- `GET /health` - 健康检查
- `GET /error` - 错误处理示例

## 技术特点

1. **独立 Go 程序**: 有自己的 go.mod 和 main.go 文件
2. **引入主项目库**: 通过 replace 指令引入主项目的功能
3. **完整功能**: 包含解析器、VM、标准库等完整功能
4. **功能单一**: 直接运行 http.zy 脚本，无需命令行参数
5. **信号处理**: 支持 Ctrl+C 优雅停止服务器

## 开发说明

### Go 模块配置

```go
module github.com/php-any/origami/examples/http

go 1.24

require (
	github.com/php-any/origami v0.0.0
	github.com/go-sql-driver/mysql v1.9.3
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/php-any/origami => ../../
```

### 主程序功能

- 创建解析器和 VM 环境
- 加载所有标准库（std、php、http、system）
- 直接运行 http.zy 脚本
- 简洁的错误处理

## 注意事项

1. 确保端口 8080 没有被其他程序占用
2. 示例中的数据结构是模拟的，实际应用中应该连接数据库
3. 可以根据需要修改端口号和路由配置
4. 在生产环境中应该添加更多的安全措施和错误处理
5. 使用 `Ctrl+C` 停止服务器
