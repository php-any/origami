<?php

namespace Bootstrap\Routing;

use Net\Http\Router;

/**
 * Route Facade（类似 Laravel Illuminate\Support\Facades\Route）
 */
class Route
{
    public static function get(string $path, array $action): void
    {
        Router::get($path, $action);
    }

    public static function post(string $path, array $action): void
    {
        Router::post($path, $action);
    }

    public static function put(string $path, array $action): void
    {
        Router::put($path, $action);
    }

    public static function delete(string $path, array $action): void
    {
        Router::delete($path, $action);
    }

    public static function group(array $attributes, callable $callback): void
    {
        Router::group($attributes, $callback);
    }

    public static function getRoutes(): array
    {
        return Router::getRoutes();
    }
}
