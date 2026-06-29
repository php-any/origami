<?php
/**
 * Wails 事件系统 + 窗口操作 示例
 *
 * 演示: 事件发射/监听、窗口尺寸/位置控制、全屏切换、主题切换
 *
 * 此脚本模拟一个"计数器"应用的后端逻辑 —
 * 实际 UI 由前端 HTML/JS 通过 Wails 绑定与 Go 后端通信。
 * 这里演示 PHP 侧如何发事件给前端，以及控制窗口行为。
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Events;
use Wails\Runtime\Log;
use Wails\Runtime\Screen;
use Wails\Runtime\Environment;

// ── 1. 构造应用 ──
$options = new App([
    'Title'     => '🔔 Wails Events Demo',
    'Width'     => 640,
    'Height'    => 480,
    'MinWidth'  => 400,
    'MinHeight' => 300,
    'MaxWidth'  => 1200,
    'MaxHeight' => 800,
]);

// ── 2. Dom 就绪 — 注册事件监听并向前端发欢迎事件 ──
$options->onDomReady(function () {
    // 获取环境信息
    $env = Environment::get();
    Log::info(sprintf(
        "环境: %s / %s / %s",
        $env[0] ?? '?', // BuildType
        $env[1] ?? '?', // Platform
        $env[2] ?? '?'  // Arch
    ));

    // 获取屏幕信息
    $screens = Screen::getAll();
    $count = count($screens);
    Log::info("检测到 {$count} 个显示器");

    foreach ($screens as $i => $screen) {
        $label = ($screen[0] ?? false) ? "当前" : (($screen[1] ?? false) ? "主" : "附加");
        Log::info("  [{$i}] {$label} — {$screen[2]}x{$screen[3]}");
    }

    // 将窗口居中
    Window::center();

    // 向前端发送初始化事件
    Events::emit("app:ready", [
        'title'   => 'Wails Events Demo',
        'message' => 'PHP 后端已就绪!',
    ]);

    Log::info("已发送 app:ready 事件到前端");
});

// ── 3. 监听前端发来的事件 ──
Events::on("counter:increment", function ($data) {
    $value = $data[0] ?? 0;
    Log::info("计数器 +1 → {$value}");
    // 更新窗口标题显示计数
    Window::setTitle("🔔 计数: {$value}");
});

Events::on("counter:decrement", function ($data) {
    $value = $data[0] ?? 0;
    Log::info("计数器 -1 → {$value}");
    Window::setTitle("🔔 计数: {$value}");
});

Events::on("window:toggleMax", function () {
    Window::toggleMaximise();
    $isMax = Window::isMaximised();
    Log::info("窗口最大化: " . ($isMax ? "是" : "否"));
    Events::emit("window:stateChanged", ['maximised' => $isMax]);
});

Events::on("window:toggleFullscreen", function () {
    if (Window::isFullscreen()) {
        Window::unfullscreen();
    } else {
        Window::fullscreen();
    }
    $isFs = Window::isFullscreen();
    Log::info("全屏: " . ($isFs ? "是" : "否"));
});

Events::on("window:resize", function ($data) {
    $w = $data[0] ?? 800;
    $h = $data[1] ?? 600;
    Window::setSize($w, $h);
    Log::info("调整窗口 → {$w}x{$h}");
});

Events::on("window:move", function ($data) {
    $x = $data[0] ?? 0;
    $y = $data[1] ?? 0;
    Window::setPosition($x, $y);
    Log::info("移动窗口 → ({$x}, {$y})");
});

Events::on("theme:dark", function () {
    Window::setDarkTheme();
    Events::emit("theme:changed", ['theme' => 'dark']);
    Log::info("切换到暗色主题");
});

Events::on("theme:light", function () {
    Window::setLightTheme();
    Events::emit("theme:changed", ['theme' => 'light']);
    Log::info("切换到亮色主题");
});

// ── 4. 退出时清理 ──
$options->onShutdown(function () {
    Events::emit("app:closing");
    Log::info("应用关闭，已发送 app:closing 事件");
});

// ── 5. 启动 ──
Application::run($options);
