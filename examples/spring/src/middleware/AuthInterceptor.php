<?php

namespace Spring\Middleware;

use Net\Http\Request;
use Net\Http\Response;

/**
 * 认证拦截器 - 实现洋葱模型中间件
 *
 * 使用方式：
 * #[Middleware(AuthInterceptor::class)]
 * #[Controller]
 * class UserController { }
 */
class AuthInterceptor {

    private array $excludePaths = [
        '/api/auth/login',
        '/api/auth/register',
        '/api/hello',
        '/api/queries',
    ];

    /**
     * 洋葱模型中间件处理
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

        \Log::info("[Auth] 认证通过: " . $path);

        $next($request, $response);

        \Log::info("[Auth] 请求处理完成");
    }

    private function shouldExclude(string $path): bool {
        foreach ($this->excludePaths as $excludePath) {
            if (strpos($path, $excludePath) === 0) {
                return true;
            }
        }
        return false;
    }

    private function verifyToken(string $token): bool {
        return !empty($token) && strlen($token) > 10;
    }
}
