<?php

namespace Spring\Middleware;

/**
 * CORS 中间件示例
 * 
 * 使用方法（在 index.php 中）：
 * $server->middleware([new CorsMiddleware(), 'handle']);
 */
class CorsMiddleware {
    
    private $allowedOrigins;
    private $allowedMethods;
    private $allowedHeaders;
    
    public function __construct(
        $allowedOrigins = ['*'],
        $allowedMethods = ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
        $allowedHeaders = ['Content-Type', 'Authorization']
    ) {
        $this->allowedOrigins = $allowedOrigins;
        $this->allowedMethods = $allowedMethods;
        $this->allowedHeaders = $allowedHeaders;
    }
    
    /**
     * 处理请求
     */
    public function handle($request, $response, $next) {
        // 设置 CORS headers
        $origin = $request->header('Origin', '*');
        
        if (in_array('*', $this->allowedOrigins) || in_array($origin, $this->allowedOrigins)) {
            $response->header('Access-Control-Allow-Origin', $origin);
        }
        
        $response->header('Access-Control-Allow-Methods', implode(', ', $this->allowedMethods));
        $response->header('Access-Control-Allow-Headers', implode(', ', $this->allowedHeaders));
        $response->header('Access-Control-Max-Age', '86400'); // 24 小时
        
        // 处理 OPTIONS 预检请求
        if ($request->method() === 'OPTIONS') {
            $response->status(204)->write('');
            return;
        }
        
        // 继续处理请求
        return $next($request, $response);
    }
}
