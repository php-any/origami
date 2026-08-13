<?php

/**
 * 官方安装就绪验收：Provider 启动 + Illuminate 官方路由可见。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\TelescopeServiceProvider as PackageTelescopeServiceProvider;

if (!class_exists(PackageTelescopeServiceProvider::class)) {
    echo "FAIL: Telescope package provider missing\n";
    exit(1);
}

$needles = [
    'telescope/telescope-api/requests',
    'telescope/telescope-api/requests/{telescopeEntryId}',
    'telescope/{view?}',
];

$set = [];
foreach (illuminate_router()->getRoutes() as $route) {
    $set[] = (string) $route->uri();
}

foreach ($needles as $needle) {
    if (!in_array($needle, $set, true)) {
        echo "FAIL: official route missing {$needle}\n";
        exit(1);
    }
}

echo "PASS\n";
