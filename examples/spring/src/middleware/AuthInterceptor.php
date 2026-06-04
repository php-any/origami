<?php

namespace Spring\Middleware;

/**
 * 认证拦截器 - 实现 Spring 风格的中间件生命周期
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
     * 前置处理 - 在控制器方法执行前调用
     * 返回 false 可以中断请求
     */
    public function preHandle($request, $response) {
        $path = $request->path();

        // 检查是否需要排除
        if ($this->shouldExclude($path)) {
            return true;
        }

        // 获取 Authorization header
        $token = $request->header('Authorization', '');

        if (empty($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "未提供认证令牌",
                "data" => null
            ]);
            return false; // 中断请求
        }

        // 验证 token（简化示例）
        if (!$this->verifyToken($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "无效的认证令牌",
                "data" => null
            ]);
            return false; // 中断请求
        }

        echo "[Auth] 认证通过: " . $path . "\n";
        return true;
    }

    /**
     * 后置处理 - 在控制器方法执行后调用
     */
    public function postHandle($request, $response) {
        echo "[Auth] 后置处理完成\n";
    }

    /**
     * 完成处理 - 整个请求完成后调用
     */
    public function afterCompletion($request, $response) {
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
