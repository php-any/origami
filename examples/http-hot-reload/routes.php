<?php

use HotReload\Controller\GreetingController;
use HotReload\Service\GreetingService;

// 这里只注册路由；控制器和服务在请求内才加载、实例化。
$server->get('/', function ($request, $response) {
    $controller = new GreetingController(new GreetingService());
    $controller->home($request, $response);
});

$server->get('/api/greeting', function ($request, $response) {
    $controller = new GreetingController(new GreetingService());
    $controller->show($request, $response);
});

$server->get('/startup', function ($request, $response) {
    $response->write(startup_version());
});
