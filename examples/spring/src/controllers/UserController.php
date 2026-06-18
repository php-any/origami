<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
use Spring\Service\UserService;
use Spring\Middleware\AuthInterceptor;
use Spring\Middleware\LogInterceptor;

#[Middleware(AuthInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class UserController {

    public function __construct(
        private UserService $userService,
    ) {}

    #[GetMapping(path: "/users")]
    public function users(Request $request, Response $response): void {
        $users = $this->userService->findAll();
        $userArray = array_map(function($user) {
            return $user->toArray();
        }, $users);
        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $userArray,
            "total" => count($userArray)
        ]);
    }

    #[GetMapping(path: "/user/{id}")]
    public function user(Request $request, Response $response): void {
        $id = (int)$request->pathValue('id');
        $user = $this->userService->findById($id);
        if (!$user) {
            $response->status(404)->json([
                "code" => 404,
                "message" => "用户不存在",
                "data" => null
            ]);
            return;
        }
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $user->toArray()
        ]);
    }

    #[PostMapping(path: "/users")]
    public function createUser(Request $request, Response $response): void {
        $body = $request->body();
        if (!isset($body['name']) || !isset($body['email'])) {
            $response->status(400)->json([
                "code" => 400,
                "message" => "缺少必要参数：name 和 email",
                "data" => null
            ]);
            return;
        }
        $user = $this->userService->create($body);
        $response->status(201)->json([
            "code" => 201,
            "message" => "created",
            "data" => $user->toArray()
        ]);
    }
}
