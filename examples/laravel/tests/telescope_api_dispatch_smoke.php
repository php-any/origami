<?php

/**
 * 官方 API 路由与 Controller 动作验收（不依赖示例侧手搓 dispatch）。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

$foundIndex = false;
$foundShow = false;

foreach (illuminate_router()->getRoutes() as $route) {
    $uri = (string) $route->uri();
    $action = (string) $route->getActionName();
    $methods = $route->methods();

    if ($uri === 'telescope/telescope-api/logs' && in_array('POST', $methods, true) && str_ends_with($action, 'LogController@index')) {
        $foundIndex = true;
    }
    if ($uri === 'telescope/telescope-api/requests/{telescopeEntryId}' && in_array('GET', $methods, true) && str_ends_with($action, 'RequestsController@show')) {
        $foundShow = true;
    }
}

if (!$foundIndex) {
    echo "FAIL: missing official LogController@index route\n";
    exit(1);
}
if (!$foundShow) {
    echo "FAIL: missing official RequestsController@show route\n";
    exit(1);
}

echo "PASS\n";
