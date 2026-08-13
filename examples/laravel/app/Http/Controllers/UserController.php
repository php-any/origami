<?php

namespace App\Http\Controllers;

use App\Services\UserService;
use Net\Http\Response;

class UserController
{
    public function __construct(
        private UserService $userService,
    ) {}

    public function index(Response $response): void
    {
        $users = $this->userService->all();
        $list = array_map(fn ($u) => $this->userService->toArray($u), $users);
        $response->success(['list' => $list, 'total' => count($list)]);
    }

    public function show(int $id, Response $response): void
    {
        $user = $this->userService->find($id);
        if ($user === null) {
            $response->error('用户不存在', 404);
            return;
        }
        $response->success($this->userService->toArray($user));
    }
}
