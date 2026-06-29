<?php
/**
 * Wails v3 Hello World — 最简桌面应用
 *
 * 演示 Wails v3 的基本使用模式:
 *   - application.Options 配置应用名称
 *   - WebviewWindowOptions 配置窗口标题和尺寸
 *   - 生命周期回调 (onStartup / onDomReady / onShutdown / onBeforeClose)
 *   - 窗口操作 (center / setTitle)
 *
 * 对应 Go API:
 *   app := application.New(application.Options{Name: "..."})
 *   app.Window.NewWithOptions(application.WebviewWindowOptions{...})
 *   app.Run()
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Log;

// ── 1. 创建应用选项 ──
// App 同时承载 application.Options 和 WebviewWindowOptions 的属性
$options = new App([
    'Title'  => '🎉 Hello Wails v3 from Origami!',
    'Width'  => 800,
    'Height' => 600,
]);

// ── 2. 生命周期回调 ──

// 应用启动时
$options->onStartup(function () {
    Log::info("Application started!");
});

// DOM 就绪时：居中窗口 + 更新标题
$options->onDomReady(function () {
    Log::info("DOM is ready!");
    Window::center();
    Window::setTitle("✅ 欢迎使用 Origami + Wails v3");
});

// 关闭前确认 (返回 true 可阻止关闭)
$options->onBeforeClose(function () {
    Log::warning("Application about to close...");
    return false;
});

// 退出时清理
$options->onShutdown(function () {
    Log::info("Application shutdown complete");
});

// ── 3. 启动应用 (阻塞) ──
Application::run($options);
