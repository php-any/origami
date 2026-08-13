<?php

/**
 * MailWatcher 官方注册状态检查。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\Telescope;
use Laravel\Telescope\Watchers\MailWatcher;

if (!Telescope::hasWatcher(MailWatcher::class)) {
    echo "FAIL: MailWatcher not registered\n";
    exit(1);
}

echo "PASS\n";
