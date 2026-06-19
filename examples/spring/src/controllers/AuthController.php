<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
use Spring\Service\AuthService;
use Spring\Middleware\LogInterceptor;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class AuthController {

    public function __construct(
        private AuthService $authService,
    ) {}

    #[PostMapping(path: "/auth/login")]
    public function login(Request $request, Response $response): void {
        $body = $request->body();

        if (!isset($body['username']) || !isset($body['password'])) {
            $response->error('缺少用户名或密码', 400);
            return;
        }

        $result = $this->authService->login($body['username'], $body['password']);

        if (!$result['success']) {
            $response->error($result['message'], 401);
            return;
        }

        $response->success([
            'token' => $result['token'],
            'user' => $result['user'],
        ], '登录成功');
    }

    #[PostMapping(path: "/auth/register")]
    public function register(Request $request, Response $response): void {
        $body = $request->body();

        if (!isset($body['username']) || !isset($body['password']) || !isset($body['email'])) {
            $response->error('缺少必要参数：username, password, email', 400);
            return;
        }

        $result = $this->authService->register($body);

        if (!$result['success']) {
            $response->error($result['message'], 400);
            return;
        }

        $response->success($result['user'], '注册成功', 201);
    }

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

    #[PostMapping(path: "/auth/logout")]
    public function logout(Request $request, Response $response): void {
        $response->success(null, '退出登录成功');
    }
}
