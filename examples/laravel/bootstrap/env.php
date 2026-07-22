<?php

/**
 * 加载 .env 到进程环境（类似 Laravel vlucas/phpdotenv 的简化版）
 */
if (!function_exists('load_env')) {
    function load_env(?string $path = null): void
    {
        if ($path === null) {
            $path = dirname(__DIR__) . '/.env';
        }

        if (!is_file($path)) {
            return;
        }

        $lines = explode("\n", file_get_contents($path));
        foreach ($lines as $line) {
            $line = trim($line);
            if ($line === '' || str_starts_with($line, '#')) {
                continue;
            }

            $pos = strpos($line, '=');
            if ($pos === false) {
                continue;
            }

            $name = trim(substr($line, 0, $pos));
            $value = trim(substr($line, $pos + 1));
            if ($value !== '' && ($value[0] === '"' || $value[0] === "'")) {
                $value = trim($value, "\"'");
            }

            if ($name === '') {
                continue;
            }

            if (getenv($name) === false) {
                putenv($name . '=' . $value);
            }
            $_ENV[$name] = $value;
            $_SERVER[$name] = $value;
        }
    }
}

if (!function_exists('env')) {
    function env(string $key, mixed $default = null): mixed
    {
        $value = getenv($key);
        if ($value === false) {
            return $default;
        }

        return match (strtolower($value)) {
            'true', '(true)' => true,
            'false', '(false)' => false,
            'null', '(null)' => null,
            default => $value,
        };
    }
}
