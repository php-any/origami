<?php

namespace Bootstrap\View;

use Net\Http\Response;

/**
 * 视图：Illuminate FileViewFinder 解析路径，Origami response->view() 渲染 HTML。
 */
class View
{
    private static string $basePath = '';

    public static function setBasePath(string $path): void
    {
        self::$basePath = rtrim($path, '/');

        $config = illuminate_config();
        $paths = $config->get('view.paths', []);
        if (!is_array($paths)) {
            $paths = [];
        }
        if (!in_array(self::$basePath, $paths, true)) {
            array_unshift($paths, self::$basePath);
        }
        $config->set('view.paths', $paths);

        $finder = illuminate_container()->make('view.finder');
        if (method_exists($finder, 'setPaths')) {
            $finder->setPaths($paths);
        }
    }

    public static function path(string $name): string
    {
        return illuminate_view()->getFinder()->find($name);
    }

    public static function render(Response $response, string $name, array $data, string $layout = 'layouts.app'): void
    {
        $response->view(self::path($name), $data, self::path($layout));
    }
}
