<?php
/**
 * Wails v3 平台专属选项 示例
 *
 * 演示 Wails v3 的:
 *   - 平台专属窗口选项 (Windows / Mac / Linux)
 *   - macOS 沉浸式标题栏 (透明 + 全尺寸内容)
 *   - 背景色 (BackgroundColour) 与 Windows 材质 (BackdropType)
 *   - 窗口启动状态 (WindowStartState)
 *   - 单实例锁 (SingleInstanceLock)
 *
 * 对应 Go API:
 *   application.WebviewWindowOptions{
 *       Windows: application.WindowsWindow{...},
 *       Mac:     application.MacWindow{...},
 *       Linux:   application.LinuxWindow{...},
 *   }
 *
 * 运行:
 *   go run ./cmd examples/platform_options.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Options\RGBA;
use Wails\Options\Windows;
use Wails\Options\Mac;
use Wails\Options\Linux;
use Wails\Options\Mac\TitleBar as MacTitleBar;
use Wails\Options\SingleInstanceLock;
use Wails\WindowStartState;
use Wails\BackdropType;
use Wails\Theme as WindowsTheme;
use Wails\WebviewGpuPolicy;
use Wails\MacAppearance;
use Wails\Runtime\Window;
use Wails\Runtime\Log;

// ── 1. 前端页面 ──
$html = <<<HTML
<!doctype html>
<html lang="zh">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Wails v3 Platform Options</title>
<style>
  :root { color-scheme: dark; }
  body {
    margin: 0; min-height: 100vh; display: flex; flex-direction: column;
    align-items: center; justify-content: center; gap: 18px;
    padding: 48px 40px 40px; /* 透明标题栏 + FullSizeContent 时为红绿灯留空 */
    font-family: -apple-system, "Segoe UI", system-ui, sans-serif; text-align: center;
    background: rgba(15, 23, 42, 0.78); color: #e2e8f0;
    --wails-draggable: drag;
  }
  h1 { margin: 0; font-size: 26px; }
  .platform { font-size: 64px; }
  .tags { display: flex; gap: 10px; flex-wrap: wrap; justify-content: center; max-width: 560px; }
  .tag { padding: 8px 14px; border-radius: 999px; background: rgba(255,255,255,.08); font-size: 13px; }
  p { opacity: .7; max-width: 520px; line-height: 1.6; margin: 0; }
</style>
</head>
<body>
  <div class="platform">🖥️</div>
  <h1>平台专属选项 Demo</h1>
  <p>本窗口根据当前操作系统应用了不同的原生选项 (标题栏、材质、GPU 策略等)。</p>
  <div class="tags">
    <span class="tag">深色背景 RGBA(30,30,30)</span>
    <span class="tag">macOS 透明标题栏</span>
    <span class="tag">Windows Mica 材质</span>
    <span class="tag">Linux GPU OnDemand</span>
    <span class="tag">单实例锁</span>
  </div>
</body>
</html>
HTML;

// ── 2. 平台专属选项 ──

// Windows: application.WindowsWindow
$winOptions = new Windows([
    'Theme'             => WindowsTheme::SYSTEM_DEFAULT,
    'BackdropType'      => BackdropType::MICA,
    'DisableWindowIcon' => false,
    'ResizeDebounceMS'  => 100,
]);

// macOS: 沉浸式标题栏 + 外观
$titleBar = new MacTitleBar([
    'TitlebarAppearsTransparent' => true,
    'FullSizeContent'            => true,
    'HideToolbarSeparator'       => true,
]);

$macOptions = new Mac([
    'Appearance' => MacAppearance::DARK_AQUA,
    'TitleBar'   => $titleBar,
]);

// Linux: application.LinuxWindow
$linuxOptions = new Linux([
    'WindowIsTranslucent' => false,
    'WebviewGpuPolicy'    => WebviewGpuPolicy::ON_DEMAND,
    'ProgramName'         => 'origami-wails-v3-demo',
]);

// 单实例锁
$singleLock = new SingleInstanceLock('com.origami.wails-v3-demo');

// ── 3. 组装应用 ──
$options = new App([
    'Title'            => '✨ Wails v3 Platform Demo',
    'Width'            => 900,
    'Height'           => 620,
    'WindowStartState' => WindowStartState::NORMAL,
    'BackgroundColour' => new RGBA(30, 30, 30, 255),
    'HTML'             => $html,

    'Windows'            => $winOptions,
    'Mac'                => $macOptions,
    'Linux'              => $linuxOptions,
    'SingleInstanceLock' => $singleLock,
]);

// ── 4. 生命周期 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("平台选项 demo 已启动, 当前 OS: " . PHP_OS);
});

$options->onShutdown(function () {
    Log::info("正在关闭...");
});

// ── 5. 启动 ──
Application::run($options);
