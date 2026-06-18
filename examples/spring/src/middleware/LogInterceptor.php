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

    /**
     * 洋葱模型中间件处理
     */
    public function handle(Request $request, Response $response, callable $next): void {
        $method = $request->method();
        $path = $request->path();

        // 前置处理
        echo "[LOG] >>> {$method} {$path}\n";

        // 调用下一个中间件或控制器
        $next($request, $response);

        // 后置处理（洋葱回溯阶段）
        echo "[LOG] <<< 响应完成\n";
        echo "[LOG] === 请求处理结束\n";
    }
}
