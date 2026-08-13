<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Operation;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
use Spring\DTO\Request\LoginRequest;
use Spring\DTO\Request\RegisterRequest;
use Spring\Service\AuthService;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class AuthController {

    public function __construct(
        private AuthService $authService,
    ) {}

    #[Operation(summary: "用户登录", description: "用户登录并返回认证令牌", tags: ["auth"])]
    #[PostMapping(path: "/auth/login")]
    public function login(LoginRequest $request, Response $response): void {
        $result = $this->authService->login($request->username, $request->password);

        if (!$result['success']) {
            $response->error($result['message'], 401);
            return;
        }

        $response->success([
            'token' => $result['token'],
            'user' => $result['user'],
        ], '登录成功');
    }

    #[Operation(summary: "用户注册", description: "用户注册并返回用户信息", tags: ["auth"])]
    #[PostMapping(path: "/auth/register")]
    public function register(RegisterRequest $request, Response $response): void {
        $result = $this->authService->register([
            'username' => $request->username,
            'password' => $request->password,
            'email' => $request->email,
        ]);

        if (!$result['success']) {
            $response->error($result['message'], 400);
            return;
        }

        $response->success($result['user'], '注册成功', 201);
    }

    #[Operation(summary: "获取用户信息", description: "获取当前登录用户信息", tags: ["user"])]
    #[GetMapping(path: "/auth/profile")]
    public function profile(Request $request, Response $response): void {
        $token = $request->header('Authorization', '');

        if (empty($token)) {
            $response->error('未提供认证令牌', 401);
            return;
        }

        $user = $this->authService->verifyToken($token);

        if (!$user) {
            $response->error('无效的认证令牌', 401);
            return;
        }

        $response->success($user);
    }

    #[Operation(summary: "用户退出登录", description: "用户退出登录并返回成功信息", tags: ["auth"])]
    #[PostMapping(path: "/auth/logout")]
    public function logout(Response $response): void {
        $response->success(null, '退出登录成功');
    }
}
