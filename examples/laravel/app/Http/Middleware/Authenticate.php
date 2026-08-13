<?php

namespace App\Http\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * API 认证中间件（Illuminate Auth Guard）
 */
class Authenticate
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        auth_set_request($request);

        if (auth()->guest()) {
            $response->error('Unauthenticated.', 401);
            return;
        }

        $next($request, $response);
    }
}
