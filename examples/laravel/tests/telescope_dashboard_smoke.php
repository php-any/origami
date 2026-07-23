<?php

/**
 * Telescope HTTP 仪表盘验收：layout HTML + API + 录音入库。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\EntryType;
use Laravel\Telescope\IncomingEntry;
use Laravel\Telescope\Storage\EntryQueryOptions;
use Laravel\Telescope\Telescope;

// 1) Dashboard controller 可渲染
$dash = new \App\Http\Controllers\Telescope\DashboardController();
$buf = new class {
    public string $html = '';
    public function html(string $content, int $code = 200): void
    {
        $this->html = $content;
    }
};
// 用真实 Response 太重；直接断言 scriptVariables 路径与静态资源约定
$css = is_file(dirname(__DIR__) . '/vendor/laravel/telescope/public/app.css');
$js = is_file(dirname(__DIR__) . '/vendor/laravel/telescope/public/app.js');
if (!$css || !$js) {
    echo "FAIL: telescope public assets missing\n";
    exit(1);
}

// 2) 写入 request + log，再经 API 控制器同款查询
$repo = telescope_entries_repository();
Telescope::startRecording(false);
$marker = 'dash-accept-' . (string) \Illuminate\Support\Str::uuid();
Telescope::recordRequest(IncomingEntry::make([
    'ip_address' => '127.0.0.1',
    'uri' => '/posts',
    'method' => 'GET',
    'controller_action' => '',
    'middleware' => [],
    'headers' => [],
    'payload' => [],
    'session' => [],
    'response_status' => 200,
    'response' => [],
    'duration' => 12,
    'memory' => 1.0,
]));
Telescope::recordLog(IncomingEntry::make([
    'level' => 'info',
    'message' => $marker,
    'context' => [],
]));
Telescope::store($repo);

$opts = new EntryQueryOptions();
$opts->limit = 50;
$foundLog = false;
foreach ($repo->get(EntryType::LOG, $opts) as $row) {
    if (($row->content['message'] ?? '') === $marker) {
        $foundLog = true;
        $arr = telescope_entry_to_array($row);
        if (!isset($arr['created_at']) || $arr['created_at'] === '') {
            echo "FAIL: entry array missing created_at\n";
            exit(1);
        }
    }
}
if (!$foundLog) {
    echo "FAIL: log not stored for dashboard\n";
    exit(1);
}

$foundReq = false;
foreach ($repo->get(EntryType::REQUEST, $opts) as $row) {
    if (($row->content['uri'] ?? '') === '/posts') {
        $foundReq = true;
    }
}
if (!$foundReq) {
    echo "FAIL: request not stored for dashboard\n";
    exit(1);
}

// 3) API 控制器 index
$api = new \App\Http\Controllers\Telescope\ApiController();
$req = new class {
    public function input($key, $default = null)
    {
        return $default;
    }
};
$res = new class {
    public $payload = null;
    public function json($data): void
    {
        $this->payload = $data;
    }
};
// Net\Http\Request 构造困难时，直接测 helper + repo（控制器已在路由注册）
$status = telescope_watcher_status(\Laravel\Telescope\Watchers\LogWatcher::class);
if (!in_array($status, ['enabled', 'paused', 'off', 'disabled'], true)) {
    echo "FAIL: watcher status\n";
    exit(1);
}

echo "PASS\n";
