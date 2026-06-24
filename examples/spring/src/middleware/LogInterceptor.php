<?php

namespace Spring\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * 日志拦截器 - 记录请求日志（洋葱模型）
 *
 * 使用方式：
 * #[Middleware(LogInterceptor::class)]
 * #[Controller]
 * class UserController { }
 */
class LogInterceptor {

    public function handle(Request $request, Response $response, callable $next): void {
        $method = $request->method();
        $path = $request->path();

        \Log::info("[LOG] >>> {$method} {$path}");

        $next($request, $response);

        \Log::info("[LOG] <<< 响应完成");
        \Log::info("[LOG] === 请求处理结束");
    }
}
