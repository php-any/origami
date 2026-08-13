<?php

/**
 * 官方 Telescope::start 后核心 watcher 已注册。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\Telescope;
use Laravel\Telescope\Watchers\CacheWatcher;
use Laravel\Telescope\Watchers\ExceptionWatcher;
use Laravel\Telescope\Watchers\LogWatcher;
use Laravel\Telescope\Watchers\QueryWatcher;
use Laravel\Telescope\Watchers\RequestWatcher;

$required = [
    RequestWatcher::class,
    LogWatcher::class,
    QueryWatcher::class,
    ExceptionWatcher::class,
    CacheWatcher::class,
];

foreach ($required as $watcher) {
    if (!Telescope::hasWatcher($watcher)) {
        echo 'FAIL: watcher not registered: ' . $watcher . "\n";
        exit(1);
    }
}

echo "PASS\n";
