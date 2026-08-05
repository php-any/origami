<?php

/**
 * LogWatcher 官方类与注册状态检查。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\Telescope;
use Laravel\Telescope\Watchers\LogWatcher;

$watcher = new LogWatcher(['enabled' => true, 'level' => 'debug']);
if (!method_exists($watcher, 'register') || !method_exists($watcher, 'recordLog')) {
    echo "FAIL: LogWatcher API missing\n";
    exit(1);
}
if (!Telescope::hasWatcher(LogWatcher::class)) {
    echo "FAIL: LogWatcher not registered\n";
    exit(1);
}

echo "PASS\n";
