<?php
/**
 * Wails 平台专属选项 + 系统托盘 示例
 *
 * 演示: Windows/macOS/Linux 平台选项、系统托盘、系统默认主题
 *
 * 运行方式:
 *   在 Go main 中加载 wails 扩展，然后执行此脚本。
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Options\RGBA;
use Wails\Options\Windows;
use Wails\Options\Mac;
use Wails\Options\Linux;
use Wails\Options\SystemTray;
use Wails\Options\Mac\TitleBar as MacTitleBar;
use Wails\Options\Mac\AboutInfo;
use Wails\Options\AssetServer;
use Wails\Options\Debug;
use Wails\Options\SingleInstanceLock;
use Wails\Options\DragAndDrop;
use Wails\WindowStartState;
use Wails\BackdropType;
use Wails\Theme as WindowsTheme;
use Wails\WebviewGpuPolicy;
use Wails\LogLevel;
use Wails\Runtime\Window;
use Wails\Runtime\Log;
use Wails\Runtime\Events;

// ── 1. 平台专属选项 ──

// Windows 选项
$winOptions = new Windows([
    'WebviewIsTransparent' => false,
    'WindowIsTranslucent'  => false,
    'Theme'                => WindowsTheme::SYSTEM_DEFAULT,
    'BackdropType'         => BackdropType::MICA,
    'DisableWindowIcon'    => false,
    'ZoomFactor'           => 1.0,
    'IsZoomControlEnabled' => true,
    'ResizeDebounceMS'     => 100,
]);

// macOS 选项
$titleBar = new MacTitleBar([
    'TitlebarAppearsTransparent' => true,
    'FullSizeContent'            => true,
    'UseToolbar'                 => false,
    'HideToolbarSeparator'       => true,
]);

$aboutInfo = new AboutInfo("My App", "Version 1.0 — Powered by Origami");

$macOptions = new Mac([
    'WebviewIsTransparent' => true,
    'WindowIsTranslucent'  => true,
    'Appearance'           => \Wails\MacAppearance::DARK_AQUA,
    'TitleBar'             => $titleBar,
    'About'                => $aboutInfo,
]);

// Linux 选项
$linuxOptions = new Linux([
    'WindowIsTranslucent' => false,
    'WebviewGpuPolicy'    => WebviewGpuPolicy::ON_DEMAND,
    'ProgramName'         => 'origami-wails-demo',
]);

// ── 2. 系统托盘 ──
$tray = new SystemTray([
    'Title'       => 'Origami Wails App',
    'Tooltip'     => '右键打开菜单',
    'StartHidden' => false,
]);

// ── 3. 调试选项 ──
$debug = new Debug([
    'OpenInspectorOnStartup' => false,
]);

// ── 4. 资源服务器 ──
$assetServer = new AssetServer([
    'AssetsDir' => './frontend/dist',
]);

// ── 5. 拖放支持 ──
$dragDrop = new DragAndDrop([
    'EnableFileDrop'     => true,
    'DisableWebViewDrop' => false,
]);

// ── 6. 单实例锁 ──
$singleLock = new SingleInstanceLock([
    'UniqueId' => 'com.origami.wails-demo',
]);

// ── 7. 组装应用 ──
$options = new App([
    'Title'              => '✨ Platform Options Demo',
    'Width'              => 1024,
    'Height'             => 768,
    'WindowStartState'   => WindowStartState::NORMAL,
    'BackgroundColour'   => new RGBA(30, 30, 30, 255),
    'AlwaysOnTop'        => false,
    'Frameless'          => false, // macOS 无边框 + Titlebar 透明 = 沉浸式

    // 拖放
    'CSSDragProperty'    => '--wails-draggable',
    'CSSDragValue'       => 'drag',
    'DragAndDrop'        => $dragDrop,

    // 安全
    'EnableFraudulentWebsiteDetection' => true,

    // 平台
    'Windows'            => $winOptions,
    'Mac'                => $macOptions,
    'Linux'              => $linuxOptions,

    // 资源
    'AssetServer'        => $assetServer,

    // 其他
    'Debug'              => $debug,
    'SingleInstanceLock' => $singleLock,
    'LogLevel'           => LogLevel::INFO,
]);

// ── 8. 生命周期 ──
$options->onDomReady(function () {
    Window::center();

    // 拖放事件 — 文件拖到窗口
    Events::on("wails:file-drop", function ($paths) {
        Log::info("收到文件: " . implode(", ", $paths));
    });

    Log::info("平台专属选项示例已启动");
    Log::info("运行平台: " . PHP_OS);
});

// ── 9. 启动 ──
Application::run($options);
