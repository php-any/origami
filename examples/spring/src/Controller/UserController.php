<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Middleware;
use Net\Http\Response;
use Spring\DTO\Request\CreateUserRequest;
use Spring\Service\QueryDemoService;
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
    public function users(Response $response): void {
        $users = $this->userService->findAll();
        $userArray = QueryDemoService::entitiesToArray($users);
        $response->success([
            'list' => $userArray,
            'total' => count($userArray),
        ]);
    }

    #[GetMapping(path: "/user/{id}")]
    public function user(int $id, Response $response): void {
        $user = $this->userService->findById($id);
        if (!$user) {
            $response->error('用户不存在', 404);
            return;
        }
        $response->success(QueryDemoService::entityToArray($user));
    }

    #[PostMapping(path: "/users")]
    public function createUser(CreateUserRequest $request, Response $response): void {
        $user = $this->userService->create([
            'name' => $request->name,
            'email' => $request->email,
            'age' => $request->age,
        ]);
        $response->success(QueryDemoService::entityToArray($user), 'created', 201);
    }
}
