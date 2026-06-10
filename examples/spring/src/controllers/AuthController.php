<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\Middleware;
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
    public function login($request, $response) {
        $body = $request->body();

        if (!isset($body['username']) || !isset($body['password'])) {
            $response->status(400)->json([
                "code" => 400,
                "message" => "缺少用户名或密码",
                "data" => null
            ]);
            return;
        }

        $result = $this->authService->login($body['username'], $body['password']);

        if (!$result['success']) {
            $response->status(401)->json([
                "code" => 401,
                "message" => $result['message'],
                "data" => null
            ]);
            return;
        }

        $response->json([
            "code" => 200,
            "message" => "登录成功",
            "data" => [
                "token" => $result['token'],
                "user" => $result['user']
            ]
        ]);
    }

    #[PostMapping(path: "/auth/register")]
    public function register($request, $response) {
        $body = $request->body();

        if (!isset($body['username']) || !isset($body['password']) || !isset($body['email'])) {
            $response->status(400)->json([
                "code" => 400,
                "message" => "缺少必要参数：username, password, email",
                "data" => null
            ]);
            return;
        }

        $result = $this->authService->register($body);

        if (!$result['success']) {
            $response->status(400)->json([
                "code" => 400,
                "message" => $result['message'],
                "data" => null
            ]);
            return;
        }

        $response->status(201)->json([
            "code" => 201,
            "message" => "注册成功",
            "data" => $result['user']
        ]);
    }

    #[GetMapping(path: "/auth/profile")]
    public function profile($request, $response) {
        $token = $request->header('Authorization', '');

        if (empty($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "未提供认证令牌",
                "data" => null
            ]);
            return;
        }

        $user = $this->authService->verifyToken($token);

        if (!$user) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "无效的认证令牌",
                "data" => null
            ]);
            return;
        }

        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $user
        ]);
    }

    #[PostMapping(path: "/auth/logout")]
    public function logout($request, $response) {
        $response->json([
            "code" => 200,
            "message" => "退出登录成功",
            "data" => null
        ]);
    }
}
