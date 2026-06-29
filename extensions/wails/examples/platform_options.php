<?php
/**
 * Wails v3 平台专属选项 + 拖放 示例
 *
 * 演示 Wails v3 的:
 *   - 平台专属窗口选项 (Windows / Mac / Linux)
 *   - macOS 沉浸式标题栏 (TitleBar: 透明 + 全尺寸内容)
 *   - 背景色和材质 (BackgroundColour / BackdropType)
 *   - 拖放支持 (DragAndDrop)
 *   - 窗口启动状态 (WindowStartState)
 *   - 单实例锁 (SingleInstanceLock)
 *
 * 对应 Go API:
 *   WebviewWindowOptions{
 *       Windows: application.WindowsWindow{...},
 *       Mac:     application.MacWindow{...},
 *       Linux:   application.LinuxWindow{...},
 *   }
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Options\RGBA;
use Wails\Options\Windows;
use Wails\Options\Mac;
use Wails\Options\Linux;
use Wails\Options\Mac\TitleBar as MacTitleBar;
use Wails\Options\Mac\AboutInfo;
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

// —— Windows 选项 ——
// 对应: application.WindowsWindow
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

// —— macOS 选项 ——
// 对应: application.MacWindow
$titleBar = new MacTitleBar([
    'TitlebarAppearsTransparent' => true,   // 透明标题栏
    'FullSizeContent'            => true,   // 内容延伸至标题栏区域
    'UseToolbar'                 => false,
    'HideToolbarSeparator'       => true,
]);

$aboutInfo = new AboutInfo(
    "Origami Wails v3 Demo",
    "Version 1.0 — Built with Origami + Wails v3"
);

$macOptions = new Mac([
    'WebviewIsTransparent' => true,
    'WindowIsTranslucent'  => true,
    'Appearance'           => \Wails\MacAppearance::DARK_AQUA,
    'TitleBar'             => $titleBar,
    'About'                => $aboutInfo,
]);

// —— Linux 选项 ——
// 对应: application.LinuxWindow
$linuxOptions = new Linux([
    'WindowIsTranslucent' => false,
    'WebviewGpuPolicy'    => WebviewGpuPolicy::ON_DEMAND,
    'ProgramName'         => 'origami-wails-v3-demo',
]);

// ── 2. 调试选项 ──
$debug = new Debug([
    'OpenInspectorOnStartup' => false,
]);

// ── 3. 拖放支持 ──
// 对应: WebviewWindowOptions.EnableDragAndDrop (v3.0.0-alpha.56+ 已重命名为 EnableFileDrop)
$dragDrop = new DragAndDrop([
    'EnableFileDrop'     => true,
    'DisableWebViewDrop' => false,
]);

// ── 4. 单实例锁 ──
$singleLock = new SingleInstanceLock([
    'UniqueId' => 'com.origami.wails-v3-demo',
]);

// ── 5. 组装应用 ──
$options = new App([
    'Title'            => '✨ Wails v3 Platform Demo',
    'Width'            => 1024,
    'Height'           => 768,
    'WindowStartState' => WindowStartState::NORMAL,

    // 深色背景
    'BackgroundColour' => new RGBA(30, 30, 30, 255),

    // 禁用调整大小
    'DisableResize'    => false,
    'Frameless'        => false,

    // 平台专属
    'Windows'          => $winOptions,
    'Mac'              => $macOptions,
    'Linux'            => $linuxOptions,

    // 其他
    'Debug'              => $debug,
    'DragAndDrop'        => $dragDrop,
    'SingleInstanceLock' => $singleLock,
    'LogLevel'           => LogLevel::INFO,
]);

// ── 6. 生命周期 ──
$options->onDomReady(function () {
    Window::center();

    // 拖放文件事件监听
    Events::on("wails:file-drop", function ($paths) {
        Log::info("Files dropped: " . implode(", ", $paths));
    });

    Log::info("Wails v3 platform options demo started");
    Log::info("Current OS: " . PHP_OS);
});

$options->onBeforeClose(function () {
    Log::info("Shutting down...");
    return false;
});

// ── 7. 启动 ──
Application::run($options);
