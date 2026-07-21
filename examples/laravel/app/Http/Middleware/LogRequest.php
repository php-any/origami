<?php

namespace App\Http\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * 请求日志中间件（类似 Laravel 中间件）
 */
class LogRequest
{
    public function handle(Request $request, Response $response, callable $next): void
    {
        $start = microtime(true);
        \Log::info('[HTTP] ' . $request->method() . ' ' . $request->path());

        $next($request, $response);

        $ms = round((microtime(true) - $start) * 1000, 2);
        \Log::info('[HTTP] 完成 ' . $request->path() . ' (' . $ms . 'ms)');
    }
}
