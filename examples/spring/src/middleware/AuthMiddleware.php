<?php

namespace Spring\Middleware;

/**
 * 认证中间件示例
 * 
 * 使用方法（在 index.php 中）：
 * $server->middleware([new AuthMiddleware(), 'handle']);
 */
class AuthMiddleware {
    
    private $excludePaths = [
        '/api/auth/login',
        '/api/auth/register',
        '/api/hello'
    ];
    
    /**
     * 处理请求
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
            return;
        }
        
        // 验证 token（简化示例）
        if (!$this->verifyToken($token)) {
            $response->status(401)->json([
                "code" => 401,
                "message" => "无效的认证令牌",
                "data" => null
            ]);
            return;
        }
        
        // Token 有效，继续处理
        Log::info("认证通过: " . $path);
        return $next($request, $response);
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
        // 这里只是简单检查 token 格式
        return !empty($token) && strlen($token) > 10;
    }
}
