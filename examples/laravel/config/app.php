<?php

/**
 * 应用配置（类似 Laravel config/app.php）
 */
return [
    'name' => env('APP_NAME', 'LaravelDemo'),
    'env' => env('APP_ENV', 'local'),
    'debug' => env('APP_DEBUG', true),
    'url' => env('APP_URL', 'http://localhost:8080'),
    'timezone' => 'Asia/Shanghai',
    'locale' => 'zh_CN',
];
