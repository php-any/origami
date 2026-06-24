<?php

namespace Spring\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * CORS 中间件示例
 *
 * 使用方法（在 index.php 中）：
 * $server->middleware(new CorsMiddleware());
 */
class CorsMiddleware {

    private array $allowedOrigins;
    private array $allowedMethods;
    private array $allowedHeaders;

    public function __construct(
        array $allowedOrigins = ['*'],
        array $allowedMethods = ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
        array $allowedHeaders = ['Content-Type', 'Authorization']
    ) {
        $this->allowedOrigins = $allowedOrigins;
        $this->allowedMethods = $allowedMethods;
        $this->allowedHeaders = $allowedHeaders;
    }

    public function handle(Request $request, Response $response, callable $next): void {
        $origin = $request->header('Origin', '*');

        if (in_array('*', $this->allowedOrigins) || in_array($origin, $this->allowedOrigins)) {
            $response->header('Access-Control-Allow-Origin', $origin);
        }

        $response->header('Access-Control-Allow-Methods', implode(', ', $this->allowedMethods));
        $response->header('Access-Control-Allow-Headers', implode(', ', $this->allowedHeaders));
        $response->header('Access-Control-Max-Age', '86400');

        if ($request->method() === 'OPTIONS') {
            $response->status(204)->write('');
            return;
        }

        $next($request, $response);
    }
}
