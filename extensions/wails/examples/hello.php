<?php
/**
 * Wails v3 Hello World — 最简桌面应用
 *
 * 演示 Wails v3 的基本使用模式:
 *   - App 配置应用名称、窗口标题和尺寸
 *   - 通过 HTML 选项提供前端页面 (内联 HTML，自动注入 wails 运行时)
 *   - 生命周期回调 (onStartup / onDomReady / onShutdown / onBeforeClose)
 *   - 窗口操作 (center / setTitle)
 *
 * 对应 Go API:
 *   app := application.New(application.Options{Name: "...", Assets: ...})
 *   app.Window.NewWithOptions(application.WebviewWindowOptions{...})
 *   app.Run()
 *
 * 运行:
 *   go run ./cmd examples/hello.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Log;

// ── 1. 前端页面 (内联 HTML) ──
$html = <<<HTML
<!doctype html>
<html lang="zh">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Hello Wails v3</title>
<style>
  :root { color-scheme: light dark; }
  * { box-sizing: border-box; }
  body {
    margin: 0; min-height: 100vh; display: flex; flex-direction: column;
    align-items: center; justify-content: center; gap: 24px;
    font-family: -apple-system, "Segoe UI", system-ui, sans-serif;
    background: linear-gradient(135deg, #6366f1, #8b5cf6, #ec4899);
    color: #fff; text-align: center;
  }
  .logo { font-size: 80px; animation: float 3s ease-in-out infinite; }
  @keyframes float { 0%,100%{ transform: translateY(0) } 50%{ transform: translateY(-14px) } }
  h1 { margin: 0; font-size: 32px; }
  p  { margin: 0; opacity: .85; }
  .count {
    font-size: 18px; padding: 12px 24px; border-radius: 999px;
    background: rgba(255,255,255,.18); backdrop-filter: blur(6px);
  }
  button {
    font: inherit; padding: 12px 28px; border: 0; border-radius: 12px; cursor: pointer;
    background: #fff; color: #6d28d9; font-weight: 600; box-shadow: 0 8px 24px rgba(0,0,0,.25);
    transition: transform .12s ease;
  }
  button:active { transform: scale(.95); }
</style>
</head>
<body>
  <div class="logo">🎉</div>
  <h1>Hello, Origami + Wails v3!</h1>
  <p>这是一个用 Origami PHP 编写的原生桌面应用。</p>
  <div class="count">点击次数: <b id="count">0</b></div>
  <button id="btn">点我 +1</button>

  <script type="module">
    import { Events } from "/wails/runtime.js";
    let n = 0;
    const el = document.getElementById('count');
    document.getElementById('btn').addEventListener('click', () => {
      el.textContent = ++n;
      // 把点击次数发送给 PHP 后端
      Events.Emit('button:clicked', n);
    });
  </script>
</body>
</html>
HTML;

// ── 2. 创建应用选项 ──
$options = new App([
    'Title'  => '🎉 Hello Wails v3 from Origami!',
    'Width'  => 800,
    'Height' => 600,
    'HTML'   => $html,
]);

// ── 3. 生命周期回调 ──

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

// 关闭前 (返回 true 可阻止关闭)
$options->onBeforeClose(function () {
    Log::warning("Application about to close...");
    return false;
});

// 退出时清理
$options->onShutdown(function () {
    Log::info("Application shutdown complete");
});

// ── 4. 监听来自前端的事件 ──
\Wails\Runtime\Events::on('button:clicked', function ($count) {
    Log::info("Button clicked {$count} time(s) in the frontend");
});

// ── 5. 启动应用 (阻塞) ──
Application::run($options);
