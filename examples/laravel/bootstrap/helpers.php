<?php

/**
 * 容器辅助函数（CLI 使用 getInstance 容器）
 */
function app_make(string $abstract): mixed
{
    return \Container\Container::getInstance()->make($abstract);
}

function resource_path(string $path = ''): string
{
    $base = dirname(__DIR__) . '/resources';
    if ($path === '') {
        return $base;
    }

    return $base . '/' . ltrim($path, '/');
}

function base_path(string $path = ''): string
{
    $base = dirname(__DIR__);
    if ($path === '') {
        return $base;
    }

    return $base . '/' . ltrim($path, '/');
}
