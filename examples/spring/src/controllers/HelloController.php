<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class HelloController {
    public int $count = 0;

    #[GetMapping(path: "/hello")]
    public function hello(Request $request, Response $response): void {
        $response->success([
            "greeting" => "Hello World!",
            "app_name" => 'Spring Demo',
            "app_version" => '1.0.0',
            "timestamp" => time(),
            "count" => $this->count++
        ]);
    }

    #[GetMapping(path: "/info")]
    public function info(Request $request, Response $response): void {
        $response->success([
            "name" => 'Spring Demo',
            "version" => '1.0.0',
            "timezone" => 'Asia/Shanghai',
            "api_prefix" => '/api',
            "api_version" => 'v1'
        ]);
    }

    #[GetMapping(path: "/status")]
    public function status(Request $request, Response $response): void {
        $response->success([
            "status" => "running",
            "uptime" => "ok",
            "timestamp" => time()
        ]);
    }
}
