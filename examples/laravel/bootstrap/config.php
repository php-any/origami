<?php

/**
 * 配置加载（类似 Laravel config/ + config() 辅助函数）
 * 若 Foundation helpers 已定义 config()，则不再重定义。
 */

if (!function_exists('config_load')) {
    function config_load(): array
    {
        static $config = null;

        if ($config === null) {
            $base = dirname(__DIR__);
            $files = [
                'app',
                'auth',
                'database',
                'view',
                'telescope',
            ];
            $config = [];
            foreach ($files as $name) {
                $path = $base . '/config/' . $name . '.php';
                if (is_file($path)) {
                    $config[$name] = require $path;
                }
            }
        }

        return $config;
    }
}

if (!function_exists('config')) {
    function config(string|array|null $key = null, mixed $default = null): mixed
    {
        if (function_exists('illuminate_config')) {
            $repo = illuminate_config();
            if ($key === null) {
                return $repo;
            }
            if (is_array($key)) {
                $repo->set($key);

                return null;
            }

            return $repo->get($key, $default);
        }

        if ($key === null || is_array($key)) {
            return $default;
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
}
