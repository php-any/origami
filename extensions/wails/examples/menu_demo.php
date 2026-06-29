<?php
/**
 * Wails v3 菜单 + 对话框 + 快捷键 示例
 *
 * 演示 Wails v3 的:
 *   - 原生应用菜单栏构建 (Menu + 子菜单 + 文本/复选/单选项)
 *   - 键盘快捷键 (Keys::cmdOrCtrl / Keys::combo)
 *   - 文件对话框 (Dialog::openFile / openDirectory / saveFile)
 *   - 消息对话框 (Dialog::message: info / question)
 *   - 菜单项点击回调 → PHP 处理
 *
 * 对应 Go API:
 *   menu := application.NewMenu()
 *   fileMenu := menu.AddSubmenu("File")
 *   fileMenu.Add("Open").SetAccelerator("CmdOrCtrl+O").OnClick(func(*Context){...})
 *   app.Menu.Set(menu)   // macOS 必须用 Application Menu
 *   window.SetMenu(menu) // Windows/Linux
 *
 * 运行:
 *   go run ./cmd examples/menu_demo.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Menu\Menu;
use Wails\Menu\Keys;
use Wails\Runtime\Window;
use Wails\Runtime\Dialog;
use Wails\Runtime\Log;
use Wails\Dialog\FileFilter;
use Wails\Dialog\OpenDialogOptions;
use Wails\Dialog\SaveDialogOptions;
use Wails\Dialog\MessageDialogOptions;
use Wails\DialogType;

// ── 1. 前端页面 ──
$html = <<<HTML
<!doctype html>
<html lang="zh">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Wails v3 Menu Demo</title>
<style>
  :root { color-scheme: light dark; }
  body {
    margin: 0; min-height: 100vh; display: flex; flex-direction: column;
    align-items: center; justify-content: center; gap: 16px; padding: 32px;
    font-family: -apple-system, "Segoe UI", system-ui, sans-serif; text-align: center;
    background: #f8fafc; color: #0f172a;
  }
  @media (prefers-color-scheme: dark) { body { background:#0f172a; color:#e2e8f0; } code{color:#a5b4fc} }
  h1 { margin: 0; font-size: 26px; }
  p  { margin: 0; max-width: 540px; line-height: 1.6; opacity: .8; }
  code { background: rgba(99,102,241,.15); padding: 2px 6px; border-radius: 6px; }
  .hint { margin-top: 8px; font-size: 13px; opacity: .6; }
</style>
</head>
<body>
  <h1>📋 菜单 &amp; 对话框 Demo</h1>
  <p>请使用窗口顶部 (macOS 为屏幕顶部) 的<strong>原生菜单栏</strong>来体验功能。</p>
  <p>
    试试快捷键：<br>
    <code>Cmd/Ctrl + O</code> 打开文件 ·
    <code>Cmd/Ctrl + S</code> 保存文件 ·
    <code>Cmd/Ctrl + Q</code> 退出
  </p>
  <p class="hint">所有菜单点击都会在终端打印日志，并可能弹出原生对话框。</p>
</body>
</html>
HTML;

// ── 2. 构建菜单栏 ──
$menu = new Menu();

// ====== File 菜单 ======
$fileMenu = new Menu();

$fileMenu->addText("打开文件...", Keys::cmdOrCtrl("o"), function () {
    $result = Dialog::openFile(new OpenDialogOptions([
        'Title'   => '打开文件',
        'Filters' => [
            new FileFilter("文本文件", "*.txt;*.md"),
            new FileFilter("所有文件", "*.*"),
        ],
    ]));
    if ($result) {
        Log::info("已打开文件: " . $result);
    } else {
        Log::info("已取消打开");
    }
});

$fileMenu->addText("打开文件夹...", Keys::cmdOrCtrl("d"), function () {
    $result = Dialog::openDirectory(new OpenDialogOptions(['Title' => '选择文件夹']));
    if ($result) {
        Log::info("已打开文件夹: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("另存为...", Keys::cmdOrCtrl("s"), function () {
    $result = Dialog::saveFile(new SaveDialogOptions([
        'Title'           => '保存文件',
        'DefaultFilename' => 'untitled.txt',
        'Filters'         => [new FileFilter("文本文件", "*.txt")],
    ]));
    if ($result) {
        Log::info("已保存到: " . $result);
    }
});

$fileMenu->addSeparator();

$fileMenu->addText("退出", Keys::cmdOrCtrl("q"), function () {
    Log::info("退出应用");
    Application::quit();
});

// ====== Edit 菜单 ======
$editMenu = new Menu();
$editMenu->addText("撤销", Keys::cmdOrCtrl("z"), function () { Log::info("撤销"); });
$editMenu->addText("重做", Keys::combo("z", ["shift"]), function () { Log::info("重做"); });
$editMenu->addSeparator();
$editMenu->addText("剪切", Keys::cmdOrCtrl("x"), function () { Log::info("剪切"); });
$editMenu->addText("复制", Keys::cmdOrCtrl("c"), function () { Log::info("复制"); });
$editMenu->addText("粘贴", Keys::cmdOrCtrl("v"), function () { Log::info("粘贴"); });

// ====== View 菜单 (复选框 + 单选) ======
$viewMenu = new Menu();
$viewMenu->addCheckbox("显示工具栏", true, "", function () { Log::info("切换: 工具栏"); });
$viewMenu->addCheckbox("显示状态栏", false, "", function () { Log::info("切换: 状态栏"); });
$viewMenu->addSeparator();
$viewMenu->addRadio("小图标", false, "", function () { Log::info("视图: 小图标"); });
$viewMenu->addRadio("大图标", true, "", function () { Log::info("视图: 大图标"); });
$viewMenu->addRadio("详细信息", false, "", function () { Log::info("视图: 详细信息"); });

// ====== Help 菜单 ======
$helpMenu = new Menu();
$helpMenu->addText("关于", "", function () {
    Dialog::message(new MessageDialogOptions([
        'Type'    => DialogType::INFO,
        'Title'   => '关于',
        'Message' => "Origami + Wails v3 Demo\n用 Origami PHP 运行时构建的原生桌面应用",
        'Buttons' => ['好的'],
    ]));
});

// ====== 组装顶层菜单栏 ======
$menu->addSubMenu("文件", $fileMenu);
$menu->addSubMenu("编辑", $editMenu);
$menu->addSubMenu("视图", $viewMenu);
$menu->addSubMenu("帮助", $helpMenu);

// ── 3. 应用配置 ──
$options = new App([
    'Title'  => '📋 Wails v3 Menu + Dialog Demo',
    'Width'  => 820,
    'Height' => 560,
    'HTML'   => $html,
    'Menu'   => $menu,
]);

// ── 4. 生命周期 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("菜单 / 对话框 demo 已就绪 — 试试菜单栏和快捷键");
});

// 关闭前确认
$options->onBeforeClose(function () {
    Dialog::message(new MessageDialogOptions([
        'Type'    => DialogType::QUESTION,
        'Title'   => '确认退出',
        'Message' => '确定要退出吗？',
        'Buttons' => ['退出', '取消'],
    ]));
    return false; // 返回 true 可阻止关闭
});

// ── 5. 启动 ──
Application::run($options);
