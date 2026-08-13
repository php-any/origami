<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Spring\Middleware\LogInterceptor;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/")]
class HomeController {

    #[GetMapping(path: "/hello")]
    public function hello(
        #[NotBlank(message: "name 不能为空")]
        #[Size(min: 2, max: 64)]
        string $name,
        #[NotBlank(message: "age 不能为空")]
        #[Size(min: 1, max: 3)]
        int $age
    ): array
    {
        return [
            "code" => 200,
            "message" => "success",
            "data" => [
                "hello" => $name,
                "age" => $age,
            ],
        ];
    }
}
