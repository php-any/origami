<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Spring\Config\AppConfig;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/")]
class HomeController {

    #[GetMapping(path: "/")]
    public function index($request, $response) {
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "app" => AppConfig::get('app.name'),
                "version" => AppConfig::get('app.version'),
                "docs" => [
                    "GET /api/hello" => "Hello World",
                    "GET /api/info" => "应用信息",
                    "GET /api/status" => "服务状态",
                    "GET /api/users" => "用户列表",
                    "GET /api/products" => "商品列表",
                ],
            ],
        ]);
    }

    #[GetMapping(path: "/hello")]
    public function hello(string $name): array
    {
        return [
            "code" => 200,
            "message" => "success",
            "data" => [
                "hello" => $name,
            ],
        ];
    }
}
