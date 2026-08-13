# Wails v3 扩展示例

本目录包含 Wails v3 扩展的完整使用示例。所有示例使用 Origami PHP 语法，且
都内置 HTML 前端页面（通过 `App` 的 `HTML` 选项），运行后即可看到真实界面，
不再是 Wails 的默认空页。

## 运行方式

**必须在 `extensions/wails` 目录下构建和运行**（不要在仓库根目录执行 `go build -o wails ./cmd/`，那会编译错误的包）。

```bash
cd extensions/wails

# 推荐：使用 Makefile
make build          # 生成 ./wailsrunner
make run-chat       # 聊天应用
make run-hello      # Hello World

# 或手动构建
go build -o wailsrunner ./cmd/
./wailsrunner examples/hello.php
./wailsrunner examples/events_demo.php
./wailsrunner examples/menu_demo.php
./wailsrunner examples/platform_options.php
./wailsrunner examples/assetdir_demo.php
./wailsrunner examples/chat_demo.php
```

## 示例列表

| 文件 | 说明 | 涵盖内容 |
|------|------|----------|
| `hello.php` | ⭐ 入门 — 最小桌面应用 | `App` + 内联 `HTML`, 生命周期回调 (`onStartup` / `onDomReady` / `onShutdown` / `onBeforeClose`), 前端→后端事件 |
| `events_demo.php` | ⭐⭐ 中级 — 事件 + 窗口控制 | `Events` (on/emit) 双向通信, 后端计数器, 窗口状态控制, 屏幕/环境信息 |
| `menu_demo.php` | ⭐⭐ 中级 — 菜单 + 对话框 | 原生 `Menu` + 子菜单 + 复选/单选, 键盘快捷键, 文件/消息对话框 |
| `platform_options.php` | ⭐⭐⭐ 高级 — 平台选项 | Windows/Mac/Linux 专属选项, macOS 沉浸式标题栏, 背景色, 单实例锁 |
| `assetdir_demo.php` + `frontend/` | ⭐⭐ 中级 — 前端目录 | `AssetDir` 加载独立前端目录, 任务清单, 后端维护状态 |
| `chat_demo.php` + `chat/` | ⭐⭐⭐ 高级 — 聊天应用 | 多频道聊天 UI, 在线用户, 斜杠命令, 机器人自动回复, 完整前后端事件通信 |

## 前端内容来源

`App` 支持三种方式提供前端页面（三选一）：

| 选项 | 说明 |
|------|------|
| `'HTML' => '<!doctype html>...'` | 内联 HTML 字符串（示例均采用此方式；需在 HTML 内 `import` `/wails/runtime.js`） |
| `'AssetDir' => __DIR__ . '/frontend'` | 本地资源目录（见 `assetdir_demo.php`，可放前端构建产物） |
| `'URL' => 'https://...'` | 直接加载外部 URL |

在内联 HTML / 资源目录模式下，前端需以 **ES 模块** 方式引入 Wails 运行时
（`/wails/runtime.js` 末尾含 `export`，不能用普通 `<script src>` 加载）：

```html
<script type="module">
  import { Events } from "/wails/runtime.js";
  Events.Emit("my:event", { foo: 1 });
  Events.On("my:reply", (ev) => console.log(ev.data));
</script>
```

PHP 后端对应 `Events::on` / `Events::emit` 双向通信。

## Wails v3 API 对照

### 应用创建

| PHP | Go (Wails v3) |
|-----|---------------|
| `new App(['Title' => '...', 'Width' => 800, 'Height' => 600])` | `application.Options{Name: "..."}` + `application.WebviewWindowOptions{Title: "...", Width: 800, Height: 600}` |
| `Application::run($options)` | `app := application.New(opts)` → `app.Window.NewWithOptions(winOpts)` → `app.Run()` |
| `Application::quit()` | `app.Quit()` |

### 窗口操作

| PHP | Go (Wails v3) |
|-----|---------------|
| `Window::setTitle("...")` | `window.SetTitle("...")` |
| `Window::center()` | `window.Center()` |
| `Window::setSize(1024, 768)` | `window.SetSize(1024, 768)` |
| `Window::fullscreen()` / `Window::unfullscreen()` | `window.Fullscreen()` / `window.UnFullscreen()` |
| `Window::maximise()` / `Window::unmaximise()` | `window.Maximise()` / `window.UnMaximise()` |
| `Window::isMaximised()` / `Window::isFullscreen()` | `window.IsMaximised()` / `window.IsFullscreen()` |
| `Window::execJS("...")` | `window.ExecJS("...")` |

### 对话框

| PHP | Go (Wails v3) |
|-----|---------------|
| `Dialog::openFile($opts)` | `app.Dialog.OpenFile().SetTitle(...).AddFilter(...).Show()` |
| `Dialog::openDirectory($opts)` | `app.Dialog.OpenDirectory().SetTitle(...).Show()` |
| `Dialog::saveFile($opts)` | `app.Dialog.SaveFile().SetTitle(...).SetDefaultFilename(...).Show()` |
| `Dialog::message($opts)` | `app.Dialog.Info().SetTitle(...).SetMessage(...).Show()` |

### 事件

| PHP | Go (Wails v3) |
|-----|---------------|
| `Events::on("name", $callback)` | `app.Event.On("name", func(e *application.CustomEvent) {})` |
| `Events::emit("name", $data)` | `app.Event.EmitEvent(&application.CustomEvent{Name: "name", Data: data})` |
| `Events::off("name")` | `app.Event.Off("name")` |

> 注：脚本在 `Application::run()` 之前注册的 `Events::on` 会被排队，待应用创建后统一生效。

### 日志

| PHP | Go (Wails v3) |
|-----|---------------|
| `Log::info("msg")` | `app.Logger.Info("msg")` |
| `Log::warning("msg")` | `app.Logger.Warning("msg")` |
| `Log::error("msg")` | `app.Logger.Error("msg")` |

## Go 入口

示例通过 `extensions/wails/cmd/main.go` 运行，它会装载标准库与 Wails 扩展，
然后执行传入的 PHP 脚本：

```go
p := parser.NewParser()
vm := runtime.NewVM(p)

std.Load(vm)
php.Load(vm)
system.Load(vm)
wails.Load(vm)

vm.LoadAndRun(os.Args[1]) // 例如 examples/hello.php
```
