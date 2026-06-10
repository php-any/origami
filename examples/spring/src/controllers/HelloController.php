<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class HelloController {
    public int $count = 0;

    #[GetMapping(path: "/hello")]
    public function hello($request, $response) {
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "greeting" => "Hello World!",
                "app_name" => 'Spring Demo',
                "app_version" => '1.0.0',
                "timestamp" => time(),
                "count" => $this->count++
            ]
        ]);
    }

    #[GetMapping(path: "/info")]
    public function info($request, $response) {
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "name" => 'Spring Demo',
                "version" => '1.0.0',
                "timezone" => 'Asia/Shanghai',
                "api_prefix" => '/api',
                "api_version" => 'v1'
            ]
        ]);
    }

    #[GetMapping(path: "/status")]
    public function status($request, $response) {
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "status" => "running",
                "uptime" => "ok",
                "timestamp" => time()
            ]
        ]);
    }
}
