<?php

namespace App\Http\Controllers;

use App\Http\Requests\LoginRequest;
use App\Services\AuthService;
use Net\Http\Request;
use Net\Http\Response;

class AuthController
{
    public function __construct(
        private AuthService $authService,
    ) {}

    public function login(LoginRequest $request, Response $response): void
    {
        $result = $this->authService->attempt($request->email, $request->password);
        if ($result === null) {
            $response->error('邮箱或密码错误', 401);
            return;
        }
        $response->success($result, '登录成功');
    }

    public function me(Request $request, Response $response): void
    {
        auth_set_request($request);
        $user = auth()->user();
        if ($user === null) {
            $response->error('Unauthenticated.', 401);
            return;
        }
        $response->success(AuthService::userToArray($user));
    }
}
