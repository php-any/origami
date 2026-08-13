<?php

namespace Bootstrap\Routing;

use Illuminate\Routing\Route as IlluminateRoute;
use Net\Http\Router as OrigamiRouter;

/**
 * Route Facade：Illuminate Routing 注册 + Origami Net\Http\Router 分发桥接。
 */
class Route
{
    public static function get(string $path, array $action): void
    {
        illuminate_router()->get($path, $action);
        OrigamiRouter::get($path, $action);
    }

    public static function post(string $path, array $action): void
    {
        illuminate_router()->post($path, $action);
        OrigamiRouter::post($path, $action);
    }

    public static function put(string $path, array $action): void
    {
        illuminate_router()->put($path, $action);
        OrigamiRouter::put($path, $action);
    }

    public static function delete(string $path, array $action): void
    {
        illuminate_router()->delete($path, $action);
        OrigamiRouter::delete($path, $action);
    }

    public static function group(array $attributes, callable $callback): void
    {
        illuminate_router()->group($attributes, function () use ($attributes, $callback) {
            OrigamiRouter::group($attributes, $callback);
        });
    }

    /**
     * 将 Illuminate Router 中的控制器路由同步到 Origami Router。
     * 仅同步形如 "Class@method" 的控制器动作；Closure 路由会被跳过。
     */
    public static function syncIlluminateRoutes(string $prefix = ''): void
    {
        $prefix = ltrim($prefix, '/');
        foreach (illuminate_router()->getRoutes() as $route) {
            if (!$route instanceof IlluminateRoute) {
                continue;
            }
            $uri = ltrim((string) $route->uri(), '/');
            if ($prefix !== '' && !str_starts_with($uri, $prefix)) {
                continue;
            }

            $uses = (string) $route->getActionName();
            if ($uses === 'Closure' || !str_contains($uses, '@')) {
                continue;
            }
            [$controller, $action] = explode('@', $uses, 2);
            if ($controller === '' || $action === '') {
                continue;
            }

            $routePath = '/' . self::origamiPathFromIlluminateUri($uri);
            $actionArray = [$controller, $action];
            foreach ($route->methods() as $method) {
                if ($method === 'HEAD') {
                    continue;
                }
                switch (strtoupper($method)) {
                    case 'GET':
                        OrigamiRouter::get($routePath, $actionArray);
                        break;
                    case 'POST':
                        OrigamiRouter::post($routePath, $actionArray);
                        break;
                    case 'PUT':
                        OrigamiRouter::put($routePath, $actionArray);
                        break;
                    case 'DELETE':
                        OrigamiRouter::delete($routePath, $actionArray);
                        break;
                }
            }
        }
    }

    /**
     * @return list<array{method: string, path: string, controller: string, action: string, middleware: list<string>, name: ?string}>
     */
    public static function getRoutes(): array
    {
        $routes = [];

        foreach (illuminate_router()->getRoutes() as $route) {
            if (!$route instanceof IlluminateRoute) {
                continue;
            }

            $catalog = self::catalogEntryFromIlluminateRoute($route);
            if ($catalog !== null) {
                $routes[] = $catalog;
            }
        }

        return $routes;
    }

    /**
     * 将 Illuminate 路由 URI 转为 Origami/Go ServeMux 路径模式。
     * 例如 telescope/{view?} → telescope/{view...}
     */
    private static function origamiPathFromIlluminateUri(string $uri): string
    {
        return (string) preg_replace('/\{([^}]+)\?\}/', '{$1...}', $uri);
    }

  private static function catalogEntryFromIlluminateRoute(IlluminateRoute $route): ?array
    {
        $methods = array_values(array_filter(
            $route->methods(),
            static fn (string $method) => $method !== 'HEAD'
        ));
        if ($methods === []) {
            return null;
        }

        $uri = $route->uri();
        $path = $uri === '' ? '/' : '/' . ltrim($uri, '/');

        $controller = '';
        $action = '';
        $uses = (string) $route->getActionName();
        if ($uses !== 'Closure' && str_contains($uses, '@')) {
            [$controller, $action] = explode('@', $uses, 2);
        } elseif ($uses !== 'Closure') {
            $action = $uses;
        }

        $middleware = [];
        foreach ($route->gatherMiddleware() as $mw) {
            if (is_string($mw) && $mw !== '') {
                $middleware[] = $mw;
            }
        }

        return [
            'method' => $methods[0],
            'path' => $path,
            'controller' => $controller,
            'action' => $action,
            'middleware' => $middleware,
            'name' => $route->getName(),
        ];
    }
}
