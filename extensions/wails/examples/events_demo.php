<?php
/**
 * Wails v3 事件系统 + 窗口控制 示例
 *
 * 演示 Wails v3 的:
 *   - 双向事件通信 (前端 wails.Events.Emit/On  ↔  PHP Events::on/emit)
 *   - 窗口状态控制 (maximise / minimise / fullscreen / setSize / center)
 *   - 屏幕信息获取 (Screen::getAll)
 *   - 环境信息获取 (Environment::get)
 *
 * 对应 Go API:
 *   app.Event.On("name", func(e *application.CustomEvent) { ... })
 *   app.Event.EmitEvent(&application.CustomEvent{Name: "name", Data: ...})
 *
 * 运行:
 *   go run ./cmd examples/events_demo.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Events;
use Wails\Runtime\Log;
use Wails\Runtime\Screen;
use Wails\Runtime\Environment;

// ── 1. 前端页面 ──
$html = <<<HTML
<!doctype html>
<html lang="zh">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Wails v3 Events Demo</title>
<style>
  :root { color-scheme: dark; }
  * { box-sizing: border-box; }
  body {
    margin: 0; min-height: 100vh; padding: 28px;
    font-family: -apple-system, "Segoe UI", system-ui, sans-serif;
    background: #0f172a; color: #e2e8f0;
  }
  h1 { font-size: 22px; margin: 0 0 4px; }
  .sub { opacity: .6; margin: 0 0 20px; font-size: 13px; }
  .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
  .card { background: #1e293b; border-radius: 14px; padding: 18px; }
  .card h2 { font-size: 13px; text-transform: uppercase; letter-spacing: .08em; opacity: .6; margin: 0 0 12px; }
  .counter { font-size: 40px; font-weight: 700; text-align: center; margin: 4px 0 12px; }
  .row { display: flex; gap: 8px; flex-wrap: wrap; }
  button {
    font: inherit; flex: 1; min-width: 84px; padding: 10px 12px; border: 0; border-radius: 10px;
    cursor: pointer; background: #6366f1; color: #fff; font-weight: 600; transition: filter .12s;
  }
  button:hover { filter: brightness(1.12); }
  button.ghost { background: #334155; }
  pre { margin: 0; font-size: 12px; white-space: pre-wrap; line-height: 1.5; opacity: .85; }
  #log { height: 120px; overflow: auto; background: #0b1220; border-radius: 10px; padding: 10px; font-size: 12px; }
  .badge { display:inline-block; padding:2px 8px; border-radius:6px; background:#334155; font-size:11px; }
</style>
</head>
<body>
  <h1>🔔 Wails v3 Events Demo</h1>
  <p class="sub">前端按钮通过事件与 Origami PHP 后端通信</p>

  <div class="grid">
    <div class="card">
      <h2>计数器 (后端维护)</h2>
      <div class="counter" id="counter">0</div>
      <div class="row">
        <button onclick="emit('counter:decrement')">− 减少</button>
        <button onclick="emit('counter:increment')">+ 增加</button>
      </div>
    </div>

    <div class="card">
      <h2>窗口控制</h2>
      <div class="row">
        <button class="ghost" onclick="emit('window:center')">居中</button>
        <button class="ghost" onclick="emit('window:maximise')">最大化</button>
      </div>
      <div class="row" style="margin-top:8px">
        <button class="ghost" onclick="emit('window:fullscreen')">全屏切换</button>
        <button class="ghost" onclick="emit('window:resize', [960, 640])">960×640</button>
      </div>
    </div>

    <div class="card" style="grid-column: 1 / -1">
      <h2>系统信息 <span class="badge" id="env">…</span></h2>
      <pre id="screens">正在获取屏幕信息…</pre>
    </div>

    <div class="card" style="grid-column: 1 / -1">
      <h2>事件日志</h2>
      <div id="log"></div>
    </div>
  </div>

  <script type="module">
    import { Events } from "/wails/runtime.js";
    // type="module" 有独立作用域，内联 onclick 在全局作用域，需把 emit 挂到 window
    window.emit = emit;
    function emit(name, data) { Events.Emit(name, data ?? null); log('→ ' + name); }
    function log(msg) {
      const d = document.getElementById('log');
      d.innerHTML += msg + '<br>'; d.scrollTop = d.scrollHeight;
    }

    // 后端推送：计数器变化
    Events.On('counter:changed', (ev) => {
      document.getElementById('counter').textContent = ev.data;
      log('← counter:changed = ' + ev.data);
    });
    // 后端推送：系统信息
    Events.On('system:info', (ev) => {
      const d = ev.data || {};
      document.getElementById('env').textContent = (d.os || '?') + ' / ' + (d.arch || '?');
      document.getElementById('screens').textContent = d.screens || '(无屏幕信息)';
    });
    // 后端推送：窗口状态
    Events.On('window:state', (ev) => log('← window:state: ' + JSON.stringify(ev.data)));

    // 页面加载完成后向后端要一次系统信息
    window.addEventListener('DOMContentLoaded', () => emit('app:requestInfo'));
  </script>
</body>
</html>
HTML;

// ── 2. 应用配置 ──
$options = new App([
    'Title'     => '🔔 Wails v3 Events Demo',
    'Width'     => 720,
    'Height'    => 640,
    'MinWidth'  => 520,
    'MinHeight' => 480,
    'HTML'      => $html,
]);

// 后端维护的计数器状态
$counter = 0;

// ── 3. DOM 就绪 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("Events demo ready");
});

// ── 4. 监听前端事件 ──

Events::on("counter:increment", function () use (&$counter) {
    $counter++;
    Log::info("Counter => {$counter}");
    Events::emit("counter:changed", $counter);
});

Events::on("counter:decrement", function () use (&$counter) {
    $counter--;
    Log::info("Counter => {$counter}");
    Events::emit("counter:changed", $counter);
});

Events::on("window:center", function () {
    Window::center();
});

Events::on("window:maximise", function () {
    Window::toggleMaximise();
    Events::emit("window:state", ['maximised' => Window::isMaximised()]);
});

Events::on("window:fullscreen", function () {
    if (Window::isFullscreen()) {
        Window::unfullscreen();
    } else {
        Window::fullscreen();
    }
    Events::emit("window:state", ['fullscreen' => Window::isFullscreen()]);
});

Events::on("window:resize", function ($size) {
    $w = max(400, $size[0] ?? 800);
    $h = max(300, $size[1] ?? 600);
    Window::setSize($w, $h);
    Window::center();
    Log::info("Resized to {$w}x{$h}");
});

// 前端请求系统信息
Events::on("app:requestInfo", function () {
    $env = Environment::get();
    $screens = Screen::getAll();

    $lines = [];
    foreach ($screens as $i => $s) {
        $primary = ($s[1] ?? false) ? " (primary)" : "";
        $lines[] = "Display {$i}: {$s[2]}x{$s[3]}{$primary}";
    }

    Events::emit("system:info", [
        'os'      => $env[0] ?? '?',
        'arch'    => $env[2] ?? '?',
        'screens' => implode("\n", $lines),
    ]);
    Log::info("Sent system info to frontend");
});

// ── 5. 退出清理 ──
$options->onShutdown(function () {
    Log::info("Events demo closing");
});

// ── 6. 启动 ──
Application::run($options);
