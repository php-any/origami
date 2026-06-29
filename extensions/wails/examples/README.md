# Wails 扩展示例

本目录包含 Wails 扩展的使用示例。所有示例使用 Origami PHP 语法。

## 示例列表

| 文件 | 说明 | 难度 |
|------|------|------|
| `hello.php` | 最简桌面应用 — 窗口配置、生命周期回调 | ⭐ 入门 |
| `menu_demo.php` | 菜单栏 + 对话框 + 快捷键 + 右键菜单 | ⭐⭐ 中级 |
| `events_demo.php` | 前后端事件通信 + 窗口操作 + 屏幕信息 | ⭐⭐ 中级 |
| `platform_options.php` | 平台专属选项 + 系统托盘 + 拖放 | ⭐⭐⭐ 高级 |

## 运行方式

### 1. 创建 Go 启动文件

在项目根目录创建 `cmd/wailsapp/main.go`:

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

### 2. 运行

```bash
cd cmd/wailsapp
go run main.go
```

### 3. 开发模式

`wails dev` 需要前端源码目录。确保 `AssetServer.AssetsDir` 指向前端构建产物目录，
或使用 Go 的 `embed.FS` 嵌入前端资源。

## PHP API 快速参考

### 应用

```php
$app = new Wails\Options\App(['Title' => 'My App', 'Width' => 800, 'Height' => 600]);
$app->onDomReady(fn() => Wails\Runtime\Window::center());
$app->onBeforeClose(fn() => false); // return true to prevent close
Wails\Application::run($app);
Wails\Application::quit();
```

### 窗口

```php
Wails\Runtime\Window::setTitle("Hello");
Wails\Runtime\Window::setSize(1024, 768);
Wails\Runtime\Window::center();
Wails\Runtime\Window::maximise();
Wails\Runtime\Window::fullscreen();
Wails\Runtime\Window::execJS("alert('hi')");
[$w, $h] = Wails\Runtime\Window::getSize();
$isFull = Wails\Runtime\Window::isFullscreen();
```

### 对话框

```php
$filters = [new Wails\Dialog\FileFilter("Images", "*.jpg;*.png")];
$path = Wails\Runtime\Dialog::openFile(new Wails\Dialog\OpenDialogOptions([
    'Title' => '打开图片', 'Filters' => $filters,
]));
Wails\Runtime\Dialog::message(new Wails\Dialog\MessageDialogOptions([
    'Type' => Wails\DialogType::INFO, 'Title' => '提示', 'Message' => '操作完成',
]));
```

### 菜单

```php
$menu = new Wails\Menu\Menu();
$menu->addText("打开", Wails\Menu\Keys::cmdOrCtrl("o"), fn() => print("Open"));
$menu->addSeparator();
$menu->addCheckbox("选项", true, "", fn($d) => print($d->Checked ? "ON" : "OFF"));
$menu->addSubMenu("子菜单", $subMenu);
```

### 事件

```php
Wails\Runtime\Events::on("event:name", fn($data) => print($data[0]));
Wails\Runtime\Events::emit("event:name", "payload");
Wails\Runtime\Events::off("event:name");
```

### 其他

```php
Wails\Runtime\Log::info("message");
Wails\Runtime\Log::error("oops");
Wails\Runtime\Browser::openURL("https://example.com");
$screens = Wails\Runtime\Screen::getAll();
$env = Wails\Runtime\Environment::get(); // [BuildType, Platform, Arch]
```
