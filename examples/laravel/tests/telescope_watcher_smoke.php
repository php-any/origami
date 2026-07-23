<?php

/**
 * Telescope LogWatcher → recordLog → store 真实录音路径。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope();

use Illuminate\Log\Events\MessageLogged;
use Illuminate\Support\Str;
use Laravel\Telescope\EntryType;
use Laravel\Telescope\IncomingEntry;
use Laravel\Telescope\Storage\EntryQueryOptions;
use Laravel\Telescope\Telescope;
use Laravel\Telescope\Watchers\LogWatcher;

$repo = telescope_entries_repository();

Telescope::startRecording(false);

// 注册 filter：依赖 collect()->every->__invoke(Closure)
Telescope::filter(function ($entry) {
    return true;
});

$marker = 'watcher smoke ' . (string) Str::uuid();

Telescope::recordLog(
    IncomingEntry::make([
        'level' => 'warning',
        'message' => $marker . ' via recordLog',
        'context' => ['from' => 'smoke'],
    ])
);

$watcher = new LogWatcher([
    'enabled' => true,
    'level' => 'debug',
]);

if (!method_exists($watcher, 'recordLog') || !method_exists($watcher, 'register')) {
    echo "FAIL: LogWatcher API\n";
    exit(1);
}

$watcher->recordLog(new MessageLogged('info', $marker . ' via LogWatcher', [
    'telescope' => ['tag:smoke'],
]));

try {
    Telescope::store($repo);
} catch (\Throwable $e) {
    echo "FAIL: Telescope::store: " . $e->getMessage() . "\n";
    exit(1);
}

$options = (new EntryQueryOptions())->limit(50);
try {
    $logs = $repo->get(EntryType::LOG, $options);
} catch (\Throwable $e) {
    echo "FAIL: repo->get: " . $e->getMessage() . "\n";
    exit(1);
}

$messages = [];
foreach ($logs as $row) {
    $c = $row->content ?? null;
    if (is_array($c) && isset($c['message'])) {
        $messages[] = $c['message'];
    }
}

if (!in_array($marker . ' via recordLog', $messages, true)) {
    echo "FAIL: recordLog not stored\n";
    var_export($messages);
    echo "\n";
    exit(1);
}

if (!in_array($marker . ' via LogWatcher', $messages, true)) {
    echo "FAIL: LogWatcher not stored\n";
    var_export($messages);
    echo "\n";
    exit(1);
}

echo "PASS\n";
