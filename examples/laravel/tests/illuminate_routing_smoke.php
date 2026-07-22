<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Events\Dispatcher;
use Illuminate\Routing\Router;

$router = new Router(new Dispatcher(), new Container());

$hello = $router->get('/hello', function () {
    return 'world';
});
$hello->name('hello');

$users = $router->get('/users/{id}', function ($id) {
    return 'user:' . $id;
});
$users->name('users.show');

// Origami 路由 URI 当前为无 leading slash 的形式（如 hello）
if ($hello->uri() !== 'hello') {
    echo "FAIL: hello uri=";
    var_export($hello->uri());
    echo "\n";
    exit(1);
}

if ($hello->getName() !== 'hello') {
    echo "FAIL: hello name\n";
    exit(1);
}

if (strpos($users->uri(), '{id}') === false) {
    echo "FAIL: users.show uri=";
    var_export($users->uri());
    echo "\n";
    exit(1);
}

$routes = $router->getRoutes();
if (count($routes) < 2) {
    echo "FAIL: route count\n";
    exit(1);
}

echo "PASS\n";
