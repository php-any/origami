<?php

namespace Spring\Service;

use Container\Singleton;
use Spring\Model\User;

#[Singleton]
class AuthService {
    
    private $users = [];
    private $tokens = [];
    
    public function __construct() {
        // 初始化示例用户数据
        $this->users = [
            "admin" => [
                "username" => "admin",
                "password" => "123456",
                "email" => "admin@example.com",
                "role" => "admin"
            ],
            "user1" => [
                "username" => "user1",
                "password" => "password",
                "email" => "user1@example.com",
                "role" => "user"
            ]
        ];
    }
    
    /**
     * 用户登录
     */
    public function login($username, $password) {
        if (!isset($this->users[$username])) {
            return [
                "success" => false,
                "message" => "用户名或密码错误"
            ];
        }
        
        $user = $this->users[$username];
        
        if ($password !== $user['password']) {
            return [
                "success" => false,
                "message" => "用户名或密码错误"
            ];
        }
        
        // 生成 token（简化示例，实际应使用 JWT）
        $token = $this->generateToken($username);
        
        return [
            "success" => true,
            "message" => "登录成功",
            "token" => $token,
            "user" => [
                "username" => $user['username'],
                "email" => $user['email'],
                "role" => $user['role']
            ]
        ];
    }
    
    /**
     * 用户注册
     */
    public function register($data) {
        $username = $data['username'];
        
        if (isset($this->users[$username])) {
            return [
                "success" => false,
                "message" => "用户名已存在"
            ];
        }
        
        // 检查邮箱是否已注册
        foreach ($this->users as $user) {
            if ($user['email'] === $data['email']) {
                return [
                    "success" => false,
                    "message" => "邮箱已被注册"
                ];
            }
        }
        
        // 创建新用户
        $this->users[$username] = [
            "username" => $username,
            "password" => $data['password'],
            "email" => $data['email'],
            "role" => $data['role'] ?? 'user'
        ];
        
        return [
            "success" => true,
            "message" => "注册成功",
            "user" => [
                "username" => $username,
                "email" => $data['email'],
                "role" => $this->users[$username]['role']
            ]
        ];
    }
    
    /**
     * 验证 Token
     */
    public function verifyToken($token) {
        if (!isset($this->tokens[$token])) {
            return null;
        }
        
        $tokenData = $this->tokens[$token];
        
        // 检查 token 是否过期（简化示例）
        if (time() > $tokenData['expires_at']) {
            unset($this->tokens[$token]);
            return null;
        }
        
        $username = $tokenData['username'];
        
        if (!isset($this->users[$username])) {
            return null;
        }
        
        $user = $this->users[$username];
        
        return [
            "username" => $user['username'],
            "email" => $user['email'],
            "role" => $user['role']
        ];
    }
    
    /**
     * 生成 Token
     */
    private function generateToken($username) {
        // 简化示例：使用 base64 编码，实际应使用 JWT
        $token = base64_encode($username . ":" . time());

        // 存储 token（实际应使用 Redis 等）
        $this->tokens[$token] = [
            "username" => $username,
            "created_at" => time(),
            "expires_at" => time() + 3600 // 1 小时后过期
        ];

        return $token;
    }
    
    /**
     * 退出登录
     */
    public function logout($token) {
        if (isset($this->tokens[$token])) {
            unset($this->tokens[$token]);
            return true;
        }
        return false;
    }
}
