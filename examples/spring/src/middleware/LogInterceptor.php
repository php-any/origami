<?php

namespace Spring\Middleware;

/**
 * 日志拦截器 - 记录请求日志
 *
 * 使用方式：
 * #[Middleware(LogInterceptor::class)]
 * #[Controller]
 * class UserController { }
 */
class LogInterceptor {

    /**
     * 前置处理 - 记录请求开始
     */
    public function preHandle($request, $response) {
        $method = $request->method();
        $path = $request->path();
        echo "[LOG] >>> {$method} {$path}\n";
        return true;
    }

    /**
     * 后置处理 - 记录响应状态
     */
    public function postHandle($request, $response) {
        echo "[LOG] <<< 响应完成\n";
    }

    /**
     * 完成处理 - 记录请求结束
     */
    public function afterCompletion($request, $response) {
        echo "[LOG] === 请求处理结束\n";
    }
}
