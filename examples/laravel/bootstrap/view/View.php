<?php

namespace Bootstrap\View;

use Net\Http\Response;

/**
 * 视图路径解析与 layout 渲染
 */
class View
{
    private static string $basePath = '';

    public static function setBasePath(string $path): void
    {
        self::$basePath = rtrim($path, '/');
    }

    public static function path(string $name): string
    {
        $relative = str_replace('.', '/', $name) . '.html';

        return self::$basePath . '/' . $relative;
    }

    public static function render(Response $response, string $name, array $data, string $layout = 'layouts.app'): void
    {
        $response->view(self::path($name), $data, self::path($layout));
    }
}
