<?php

/**
 * 容器 / 路径辅助函数。
 * Foundation helpers 加载后优先委托官方实现；此处提供 fallback。
 */

if (!function_exists('app_make')) {
    function app_make(string $abstract): mixed
    {
        return \Container\Container::getInstance()->make($abstract);
    }
}

if (!function_exists('resource_path')) {
    function resource_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/resources';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('base_path')) {
    function base_path(string $path = ''): string
    {
        $base = dirname(__DIR__);
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('storage_path')) {
    function storage_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/storage';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('public_path')) {
    function public_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/public';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('config_path')) {
    function config_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/config';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('database_path')) {
    function database_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/database';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}

if (!function_exists('app_path')) {
    function app_path(string $path = ''): string
    {
        $base = dirname(__DIR__) . '/app';
        if ($path === '') {
            return $base;
        }

        return $base . '/' . ltrim($path, '/');
    }
}
