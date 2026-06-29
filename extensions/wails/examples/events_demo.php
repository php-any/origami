<?php
/**
 * Wails v3 事件系统 + 窗口操作 示例
 *
 * 演示 Wails v3 的:
 *   - 事件通信模型 (Events::on / Events::emit / Events::off)
 *   - app.OnEvent / app.EmitEvent 对应关系
 *   - 窗口状态控制 (maximise / minimise / fullscreen / setSize / setPosition)
 *   - 主题切换 (setDarkTheme / setLightTheme / setSystemDefaultTheme)
 *   - 屏幕信息获取 (Screen::getAll)
 *   - 环境信息获取 (Environment::get)
 *
 * 对应 Go API:
 *   app.OnEvent("event:name", func(e *application.CustomEvent) { ... })
 *   app.EmitEvent("event:name", data)
 *   app.OffEvent("event:name")
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
    'Title'     => '🔔 Wails v3 Events Demo',
    'Width'     => 640,
    'Height'    => 480,
    'MinWidth'  => 400,
    'MinHeight' => 300,
    'MaxWidth'  => 1200,
    'MaxHeight' => 800,
]);

// ── 2. Dom 就绪 — 初始化事件和获取系统信息 ──
$options->onDomReady(function () {
    // —— 获取环境信息 ——
    $env = Environment::get();
    Log::info(sprintf(
        "Environment: %s / %s / %s",
        $env[0] ?? '?',  // BuildType
        $env[1] ?? '?',  // Platform
        $env[2] ?? '?'   // Arch
    ));

    // —— 获取屏幕信息 ——
    $screens = Screen::getAll();
    $count = count($screens);
    Log::info("Detected {$count} display(s)");

    foreach ($screens as $i => $screen) {
        $isCurrent  = ($screen[0] ?? false) ? "current" : "";
        $isPrimary  = ($screen[1] ?? false) ? "primary" : "";
        $resolution = "{$screen[2]}x{$screen[3]}";
        Log::info("  Display {$i}: {$resolution} {$isCurrent} {$isPrimary}");
    }

    // 居中窗口
    Window::center();

    // 向前端发送 ready 事件
    Events::emit("app:ready", [
        'title'   => 'Wails v3 Events Demo',
        'message' => 'Backend is ready!',
    ]);

    Log::info("Sent 'app:ready' event to frontend");
});

// ── 3. 监听前端事件 (对应 JS: wails.Events.Emit({name: "...", data: ...})) ──

// 计数器事件
Events::on("counter:increment", function ($data) {
    $value = $data[0] ?? 0;
    Log::info("Counter ++ => {$value}");
    Window::setTitle("🔔 Counter: {$value}");
});

Events::on("counter:decrement", function ($data) {
    $value = $data[0] ?? 0;
    Log::info("Counter -- => {$value}");
    Window::setTitle("🔔 Counter: {$value}");
});

// 窗口状态控制事件
Events::on("window:toggleMax", function () {
    Window::toggleMaximise();
    $isMax = Window::isMaximised();
    Log::info("Window maximised: " . ($isMax ? "YES" : "NO"));
    Events::emit("window:stateChanged", [
        'maximised' => $isMax,
    ]);
});

Events::on("window:toggleFullscreen", function () {
    if (Window::isFullscreen()) {
        Window::unfullscreen();
    } else {
        Window::fullscreen();
    }
    $isFs = Window::isFullscreen();
    Log::info("Fullscreen: " . ($isFs ? "YES" : "NO"));
});

Events::on("window:resize", function ($data) {
    $w = max(400, $data[0] ?? 800);
    $h = max(300, $data[1] ?? 600);
    Window::setSize($w, $h);
    Log::info("Resized to {$w}x{$h}");
});

Events::on("window:move", function ($data) {
    $x = $data[0] ?? 0;
    $y = $data[1] ?? 0;
    Window::setPosition($x, $y);
    Log::info("Moved to ({$x}, {$y})");
});

// 主题切换事件
Events::on("theme:dark", function () {
    Window::setDarkTheme();
    Events::emit("theme:changed", ['theme' => 'dark']);
    Log::info("Switched to dark theme");
});

Events::on("theme:light", function () {
    Window::setLightTheme();
    Events::emit("theme:changed", ['theme' => 'light']);
    Log::info("Switched to light theme");
});

Events::on("theme:system", function () {
    Window::setSystemDefaultTheme();
    Events::emit("theme:changed", ['theme' => 'system']);
    Log::info("Switched to system default theme");
});

// ── 4. 退出时清理 ──
$options->onShutdown(function () {
    Events::emit("app:closing");
    Log::info("Application closing, sent 'app:closing' event");
});

// ── 5. 启动 ──
Application::run($options);
