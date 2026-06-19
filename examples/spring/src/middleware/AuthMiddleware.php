<?php

namespace Spring\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * 认证中间件示例
 *
 * 使用方法（在 index.php 中）：
 * $server->middleware(new AuthMiddleware());
 */
class AuthMiddleware {

    private array $excludePaths = [
        '/api/auth/login',
        '/api/auth/register',
        '/api/hello'
    ];

    /**
     * 处理请求
     */
    public function handle(Request $request, Response $response, callable $next): void {
        $path = $request->path();

        // 检查是否需要排除
        if ($this->shouldExclude($path)) {
            $next($request, $response);
            return;
        }

        // 获取 Authorization header
        $token = $request->header('Authorization', '');

        if (empty($token)) {
            $response->error('未提供认证令牌', 401);
            return;
        }

        // 验证 token（简化示例）
        if (!$this->verifyToken($token)) {
            $response->error('无效的认证令牌', 401);
            return;
        }

        // Token 有效，继续处理
        Log::info("认证通过: " . $path);
        $next($request, $response);
    }

    /**
     * 检查路径是否应该排除
     */
    private function shouldExclude(string $path): bool {
        foreach ($this->excludePaths as $excludePath) {
            if (strpos($path, $excludePath) === 0) {
                return true;
            }
        }
        return false;
    }

    /**
     * 验证 Token
     */
    private function verifyToken(string $token): bool {
        // 简化示例：实际应使用 JWT 验证
        // 这里只是简单检查 token 格式
        return !empty($token) && strlen($token) > 10;
    }
}
