<?php

namespace App\Http\Middleware;

use App\Services\AuthService;
use Net\Http\Request;
use Net\Http\Response;

/**
 * API 认证中间件（类似 Laravel auth middleware）
 */
class Authenticate
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        $user = AuthService::userFromRequest($request);

        if ($user === null) {
            $response->error('Unauthenticated.', 401);
            return;
        }

        $next($request, $response);
    }
}
