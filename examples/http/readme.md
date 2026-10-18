# 简易 HTTP 服务

用 Origami 的 `Net\Http\Server` 提供两个固定响应。脚本只在启动时执行一次并注册路由，之后每个请求由已注册的闭包处理。

不扫描注解、不挂中间件、不按请求打日志。`/json` 直接写出常量字符串，不在热路径上组数组再序列化。

在仓库根目录启动：

```bash
go run ./zy.go examples/http/index.php
```

默认监听 `0.0.0.0:8080`。`HTTP_HOST`、`HTTP_PORT` 可覆盖。

| 方法 | 路径 | 响应 |
|---|---|---|
| GET | `/` | `text/plain` `ok` |
| GET | `/json` | `application/json` `{"ok":true}` |
