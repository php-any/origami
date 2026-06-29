# Wails v3 扩展示例

本目录包含 Wails v3 扩展的完整使用示例。所有示例使用 Origami PHP 语法。

## 示例列表

| 文件 | 说明 | 涵盖内容 |
|------|------|----------|
| `hello.php` | ⭐ 入门 — 最小桌面应用 | `App`, `Window`, 生命周期回调 (`onStartup` / `onDomReady` / `onShutdown` / `onBeforeClose`) |
| `menu_demo.php` | ⭐⭐ 中级 — 菜单 + 对话框 | `Menu`, `MenuItem`, 键盘快捷键, 文件/消息对话框, 关闭确认 |
| `events_demo.php` | ⭐⭐ 中级 — 事件 + 窗口控制 | `Events` (on/emit/off), 窗口状态控制, 主题切换, 屏幕/环境信息 |
| `platform_options.php` | ⭐⭐⭐ 高级 — 平台选项 | Windows/Mac/Linux 专属选项, macOS 沉浸式标题栏, 拖放, 单实例锁 |

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
| `Window::setDarkTheme()` / `Window::setLightTheme()` | `window.SetDarkTheme()` / `window.SetLightTheme()` |
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
| `Events::on("name", $callback)` | `app.OnEvent("name", func(e *application.CustomEvent) {})` |
| `Events::emit("name", $data)` | `app.EmitEvent("name", data)` |
| `Events::off("name")` | `app.OffEvent("name")` |

### 日志

| PHP | Go (Wails v3) |
|-----|---------------|
| `Log::info("msg")` | `app.Logger.Info("msg")` |
| `Log::warning("msg")` | `app.Logger.Warning("msg")` |
| `Log::error("msg")` | `app.Logger.Error("msg")` |

## 运行方式

```bash
# 在含有 Go main 入口的项目中加载扩展:
cd cmd/wailsapp
go run main.go
```

Go 启动代码示例:

```go
package main

import (
    "github.com/php-any/origami"
    wailsExt "github.com/php-any/origami-wails"
)

func main() {
    origami.RunWith(func(vm *runtime.VM) {
        wailsExt.Load(vm)
        vm.LoadAndRun("examples/hello.php")
    })
}
```
