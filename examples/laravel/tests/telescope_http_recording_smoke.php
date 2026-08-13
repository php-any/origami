<?php

/**
 * RequestHandled / MessageLogged 事件类可加载。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Illuminate\Foundation\Http\Events\RequestHandled;
use Illuminate\Log\Events\MessageLogged;

if (!class_exists(RequestHandled::class) || !class_exists(MessageLogged::class)) {
    echo "FAIL: recording events missing\n";
    exit(1);
}

echo "PASS\n";
