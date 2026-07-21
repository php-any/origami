<?php

namespace App\Services;

use App\Models\ApiToken;
use App\Models\User;
use Container\Annotation\Singleton;
use Database\DB;

/**
 * 认证服务（类似 Laravel Auth）
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

        $apiToken = new ApiToken();
        $apiToken->user_id = $user->id;
        $apiToken->token = $token;
        DB::insert($apiToken);

        $info = $this->userService->toArray($user);

        return [
            'token' => $token,
            'user' => $info,
        ];
    }

    public static function userFromToken(string $token): ?array
    {
        if ($token === '') {
            return null;
        }

        $row = DB::model(ApiToken::class)->where('token = ?', $token)->first();
        if ($row === null) {
            return null;
        }

        $user = DB::model(User::class)->where('id = ?', $row->user_id)->first();
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

    public static function userFromRequest($request): ?array
    {
        $token = $request->header('Authorization', '');

        return self::userFromToken($token);
    }
}
