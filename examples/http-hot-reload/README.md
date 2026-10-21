# 控制器分层与 HTTP 热重载验证

通过现有 Origami CLI 运行本语言代码，使用 `Net\Http\Server`，不增加 Go Handler。

## 结构

- `index.php`：创建并启动 HTTP 服务。
- `bootstrap.php`：注册类映射 autoload，加载启动期对照文件。
- `routes.php`：注册路由，在请求回调中创建控制器并注入服务。
- `src/Controller/GreetingController.php`：读取查询参数，提供 JSON 和 HTML 响应。
- `src/Service/GreetingService.php`：封装问候语逻辑。
- `src/View/GreetingView.php`：加载视图模板。
- `views/greeting.php`：HTML 模板，转义动态内容。
- `startup.php`：启动期加载代码的对照。

## 运行

在仓库根目录执行，先确认 18087 端口空闲：

```powershell
go run -mod=mod ./zy.go examples/http-hot-reload/index.php
```

```powershell
curl.exe --max-time 10 'http://127.0.0.1:18087/api/greeting?name=Codex'
curl.exe --max-time 10 'http://127.0.0.1:18087/?name=Codex'
curl.exe --max-time 10 'http://127.0.0.1:18087/startup'
```

JSON 初始响应：

```json
{"controller":"controller-v1","message":"service-v1: Hello, Codex"}
```

## 验证方法与当前结果

保持同一服务进程运行，逐项修改并保存，再请求：

1. 控制器里的 `controller-v1` 改为 `controller-v2`。
2. 服务里的 `service-v1` 改为 `service-v2`。
3. 模板里的 `view-v1` 改为 `view-v2`。
4. `startup.php` 的 `v1` 改为 `v2`，作为启动期对照。

预期：前三项应在后续请求中更新，第四项需要重启。

2026-10-21 按 RequestVM 独占解析结果的设计调整后实测：控制器、服务、模板逐项修改后均立即生效；启动期对照保持旧版本。JSON、HTML、查询参数和 HTML 转义正常，16 个并发请求的响应未混用。测试过程中 PID 保持不变。

请求期的 AST、类声明和新增类路径映射归当前 RequestVM 所有，同一请求复用本地解析缓存；新请求重新读盘解析，不依赖文件时间戳变化。启动期已加载的定义仍从 Base VM 继承。测试不手动清缓存或重启。自动回归：`go test -mod=mod ./std/net/http -run TestControllerRequestParseLifecycle -count=1`。

当前运行时没有注册 `strncmp`，因此示例 autoload 使用明确的类映射，不依赖该函数。

验证后恢复文件；在服务终端按 Ctrl+C 停止。可通过 `HTTP_PORT` 环境变量指定其他端口。启动入口、路由表和 bootstrap 的变更需要重启。
