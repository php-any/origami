<?php

namespace App\Services;

use App\Models\ApiToken;
use App\Models\User;
use Container\Annotation\Singleton;

/**
 * 认证服务：登录发 token；当前用户由 illuminate/auth Guard 解析
 */
#[Singleton]
class AuthService
{
    private const PEPPER = 'origami-laravel-demo';

    public function __construct(
        private UserService $userService,
    ) {}

    public static function hashPassword(string $password): string
    {
        return hash('sha256', self::PEPPER . $password);
    }

    public function attempt(string $email, string $password): ?array
    {
        $user = $this->userService->findByEmail($email);
        if ($user === null) {
            return null;
        }

        if ($user->password !== self::hashPassword($password)) {
            return null;
        }

        $token = bin2hex(random_bytes(16));

        ApiToken::query()->create([
            'user_id' => $user->id,
            'token' => $token,
            'created_at' => date('Y-m-d H:i:s'),
        ]);

        return [
            'token' => $token,
            'user' => $this->userService->toArray($user),
        ];
    }

    public static function userFromToken(string $token): ?User
    {
        auth_set_token($token);
        $user = auth()->user();

        return $user instanceof User ? $user : null;
    }

    public static function userFromRequest($request): ?User
    {
        auth_set_request($request);
        $user = auth()->user();

        return $user instanceof User ? $user : null;
    }

    public static function userToArray(?User $user): ?array
    {
        if ($user === null) {
            return null;
        }

        return [
            'id' => $user->id,
            'name' => $user->name,
            'email' => $user->email,
            'created_at' => $user->created_at,
        ];
    }
}
