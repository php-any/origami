<?php
/**
 * Wails Hello World — 最简桌面应用
 *
 * 演示: 窗口配置、DomReady 回调、窗口居中
 *
 * 运行方式:
 *   在 Go main 中加载 wails 扩展，然后执行此脚本。
 *   详见 README.md 中的 Go 启动代码。
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Log;

// ── 1. 配置应用选项 ──
$options = new App([
    'Title'  => '🎉 Hello Wails from Origami!',
    'Width'  => 800,
    'Height' => 600,
    'DisableResize' => false,
    'Frameless'      => false,
    // 背景色 (可选)
    // 'BackgroundColour' => new \Wails\Options\RGBA(240, 248, 255, 255),
]);

// ── 2. Dom 就绪时居中窗口 ──
$options->onDomReady(function () {
    Log::info("应用已启动 — DOM 就绪");
    Window::center();
    Window::setTitle("✅ 欢迎使用 Origami + Wails");
});

// ── 3. 关闭前确认 ──
$options->onBeforeClose(function () {
    Log::warning("应用即将关闭...");
    return false; // 返回 true 可阻止关闭
});

// ── 4. 退出时清理 ──
$options->onShutdown(function () {
    Log::info("应用已退出");
});

// ── 5. 启动 ──
Application::run($options);
