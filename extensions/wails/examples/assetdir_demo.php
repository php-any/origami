<?php
/**
 * Wails v3 前端目录 (AssetDir) 示例 — 任务清单
 *
 * 与其它示例不同，本示例的前端不是内联 HTML 字符串，而是一个独立目录：
 *
 *   examples/frontend/
 *     ├── index.html   入口页面
 *     ├── style.css    样式
 *     └── app.js       逻辑 (ES 模块, import "/wails/runtime.js")
 *
 * 通过 App 的 'AssetDir' 选项指定该目录，Wails 会把目录内文件作为静态资源服务，
 * 前端再通过事件与本 PHP 后端通信。后端在内存中维护任务列表状态。
 *
 * 对应 Go API:
 *   application.Options{
 *       Assets: application.AssetOptions{
 *           Handler: application.BundledAssetFileServer(os.DirFS("frontend")),
 *       },
 *   }
 *
 * 运行:
 *   go run ./cmd examples/assetdir_demo.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Events;
use Wails\Runtime\Log;

// ── 1. 应用配置：前端来自目录而非内联 HTML ──
$options = new App([
    'Title'     => '📝 Wails v3 · 前端目录 (AssetDir)',
    'Width'     => 560,
    'Height'    => 640,
    'MinWidth'  => 420,
    'MinHeight' => 480,
    'AssetDir'  => __DIR__ . '/frontend',
]);

// ── 2. 后端状态：内存中的任务列表 ──
$todos  = [];   // id => ['id' => int, 'text' => string, 'done' => bool]
$nextId = 1;

// 把最新列表推送给前端
$broadcast = function () use (&$todos) {
    Events::emit("todo:changed", array_values($todos));
};

// ── 3. 事件处理 ──

// 前端请求当前列表 (页面加载时)
Events::on("todo:list", function () use ($broadcast) {
    $broadcast();
});

// 添加任务
Events::on("todo:add", function ($text) use (&$todos, &$nextId, $broadcast) {
    $text = trim((string)$text);
    if ($text === "") {
        return;
    }
    $id = $nextId++;
    $todos[$id] = ['id' => $id, 'text' => $text, 'done' => false];
    Log::info("添加任务 #{$id}: {$text}");
    $broadcast();
});

// 切换完成状态
Events::on("todo:toggle", function ($id) use (&$todos, $broadcast) {
    $id = (int)$id;
    if (isset($todos[$id])) {
        $todos[$id]['done'] = !$todos[$id]['done'];
        $broadcast();
    }
});

// 删除任务
Events::on("todo:delete", function ($id) use (&$todos, $broadcast) {
    $id = (int)$id;
    if (isset($todos[$id])) {
        unset($todos[$id]);
        Log::info("删除任务 #{$id}");
        $broadcast();
    }
});

// 清除所有已完成
Events::on("todo:clear", function () use (&$todos, $broadcast) {
    foreach ($todos as $id => $t) {
        if ($t['done']) {
            unset($todos[$id]);
        }
    }
    Log::info("已清除完成的任务");
    $broadcast();
});

// ── 4. 生命周期 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("AssetDir demo 就绪 — 前端来自 examples/frontend 目录");
});

// ── 5. 启动 ──
Application::run($options);
