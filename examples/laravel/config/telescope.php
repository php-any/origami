<?php

/**
 * Telescope 配置（冒烟验收用；storage 默认 sqlite）
 */

use Laravel\Telescope\Watchers;

return [
    'domain' => env('TELESCOPE_DOMAIN'),
    'path' => env('TELESCOPE_PATH', 'telescope'),
    'driver' => env('TELESCOPE_DRIVER', 'database'),
    'storage' => [
        'database' => [
            'connection' => env('DB_CONNECTION', 'sqlite'),
            'chunk' => 1000,
        ],
    ],
    'enabled' => env('TELESCOPE_ENABLED', true),
    // Origami HTTP 不走官方 Authorize；鉴权在 bootstrap_telescope_http 放行
    'middleware' => [],
    'only_paths' => [],
    'ignore_paths' => [
        'telescope*',
        'vendor/telescope*',
    ],
    'watchers' => [
        Watchers\RequestWatcher::class => [
            'enabled' => env('TELESCOPE_REQUEST_WATCHER', true),
            'size_limit' => 64,
        ],
        Watchers\LogWatcher::class => [
            'enabled' => env('TELESCOPE_LOG_WATCHER', true),
            'level' => 'debug',
        ],
        Watchers\QueryWatcher::class => [
            'enabled' => env('TELESCOPE_QUERY_WATCHER', true),
            'ignore_packages' => true,
            'ignore_paths' => [],
            'slow' => 100,
        ],
    ],
];
