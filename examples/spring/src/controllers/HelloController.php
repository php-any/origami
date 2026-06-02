<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Spring\Config\AppConfig;

#[Controller]
#[Route(prefix: "/api")]
class HelloController {
    
    #[GetMapping(path: "/hello")]
    public function hello($request, $response) {
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "greeting" => "Hello World!",
                "app_name" => AppConfig::get('app.name'),
                "app_version" => AppConfig::get('app.version'),
                "timestamp" => time()
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
                "name" => AppConfig::get('app.name'),
                "version" => AppConfig::get('app.version'),
                "timezone" => AppConfig::get('app.timezone'),
                "api_prefix" => AppConfig::get('api.prefix'),
                "api_version" => AppConfig::get('api.version')
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


