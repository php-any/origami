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
        $response->success([
            'list' => $userArray,
            'total' => count($userArray),
        ]);
    }

    #[GetMapping(path: "/user/{id}")]
    public function user(Request $request, Response $response): void {
        $id = (int)$request->pathValue('id');
        $user = $this->userService->findById($id);
        if (!$user) {
            $response->error('用户不存在', 404);
            return;
        }
        $response->success($user->toArray());
    }

    #[PostMapping(path: "/users")]
    public function createUser(Request $request, Response $response): void {
        $body = $request->body();
        if (!isset($body['name']) || !isset($body['email'])) {
            $response->error('缺少必要参数：name 和 email', 400);
            return;
        }
        $user = $this->userService->create($body);
        $response->success($user->toArray(), 'created', 201);
    }
}
