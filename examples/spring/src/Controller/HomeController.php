<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/")]
class HomeController {

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
