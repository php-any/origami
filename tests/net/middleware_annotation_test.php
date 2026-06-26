<?php

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Net\Annotation\Route;

// 测试 Spring 风格的中间件注解

// 定义认证中间件类
class AuthMiddleware {
    public function preHandle($request, $response) {
        $token = $request->Header->Get("Authorization");
        if ($token == "") {
            echo "[AUTH] 未授权访问，请求被拒绝\n";
            return false; // 中断请求
        }
        echo "[AUTH] 认证通过\n";
        return true;
    }

    public function postHandle($request, $response) {
        echo "[AUTH] 后置处理\n";
    }

    public function afterCompletion($request, $response) {
        echo "[AUTH] 完成处理\n";
    }
}

// 定义日志中间件类
class LogMiddleware {
    public function preHandle($request, $response) {
        echo "[LOG] 请求路径: " . $request->URL()->Path . "\n";
        return true;
    }

    public function postHandle($request, $response) {
        echo "[LOG] 响应完成\n";
    }

    public function afterCompletion($request, $response) {
        echo "[LOG] 请求结束\n";
    }
}

// 使用中间件注解
#[Middleware(AuthMiddleware::class)]
#[Controller("user")]
#[Route("/api")]
class UserController {

    // 方法级别中间件（会与类级别中间件合并）
    #[Middleware(LogMiddleware::class)]
    #[GetMapping("/profile")]
    public function getProfile($r, $w) {
        echo "获取用户资料\n";
    }

    // 只有类级别中间件
    #[GetMapping("/list")]
    public function getList($r, $w) {
        echo "获取用户列表\n";
    }
}

echo "测试 @Middleware 注解定义成功\n";
