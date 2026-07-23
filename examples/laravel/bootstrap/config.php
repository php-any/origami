<?php

/**
 * 配置加载（类似 Laravel config/ + config() 辅助函数）
 */

function config_load(): array
{
    static $config = null;

    if ($config === null) {
        $base = dirname(__DIR__);
        $config = [
            'app' => require $base . '/config/app.php',
            'auth' => require $base . '/config/auth.php',
            'database' => require $base . '/config/database.php',
            'view' => require $base . '/config/view.php',
        ];
    }

    return $config;
}

function config(string $key, mixed $default = null): mixed
{
    if (function_exists('illuminate_config')) {
        return illuminate_config()->get($key, $default);
    }

    $segments = explode('.', $key);
    $value = config_load();

    foreach ($segments as $segment) {
        if (!is_array($value) || !array_key_exists($segment, $value)) {
            return $default;
        }
        $value = $value[$segment];
    }

    return $value;
}
