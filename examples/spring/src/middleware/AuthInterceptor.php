<?php

namespace Spring\Middleware;

/**
 * 认证拦截器 - 实现洋葱模型中间件
 *
 * 使用方式：
 * #[Middleware(AuthInterceptor::class)]
 * #[Controller]
 * class UserController { }
 */
class AuthInterceptor {

    private $excludePaths = [
        '/api/auth/login',
        '/api/auth/register',
        '/api/hello'
    ];

    /**
     * 洋葱模型中间件处理
     *
     * @param $request  请求对象
     * @param $response 响应对象
     * @param $next     调用下一个中间件或控制器的回调
     */
    public function handle($request, $response, $next) {
        $path = $request->path();

        // 检查是否需要排除
        if ($this->shouldExclude($path)) {
            return $next($request, $response);
        }

        // 获取 Authorization header
        $token = $request->header('Authorization', '');

        if (empty($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "未提供认证令牌",
                "data" => null
            ]);
            return; // 不调用 $next，中断请求
        }

        // 验证 token（简化示例）
        if (!$this->verifyToken($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "无效的认证令牌",
                "data" => null
            ]);
            return; // 不调用 $next，中断请求
        }

        echo "[Auth] 认证通过: " . $path . "\n";

        // 调用下一个中间件或控制器
        $next($request, $response);

        // 后置处理（洋葱回溯阶段）
        echo "[Auth] 请求处理完成\n";
    }

    /**
     * 检查路径是否应该排除
     */
    private function shouldExclude($path) {
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
    private function verifyToken($token) {
        // 简化示例：实际应使用 JWT 验证
        return !empty($token) && strlen($token) > 10;
    }
}
