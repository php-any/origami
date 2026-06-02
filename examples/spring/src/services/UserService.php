<?php

namespace Spring\Service;

use Spring\Model\User;

class UserService {
    
    private $users = [];
    
    public function __construct() {
        // 初始化示例数据
        $this->users = [
            1 => new User(1, "张三", "zhangsan@example.com", 25),
            2 => new User(2, "李四", "lisi@example.com", 30),
            3 => new User(3, "王五", "wangwu@example.com", 28)
        ];
    }
    
    /**
     * 获取所有用户
     */
    public function findAll() {
        return array_values($this->users);
    }
    
    /**
     * 根据 ID 查找用户
     */
    public function findById($id) {
        return $this->users[$id] ?? null;
    }
    
    /**
     * 根据邮箱查找用户
     */
    public function findByEmail($email) {
        foreach ($this->users as $user) {
            if ($user->getEmail() === $email) {
                return $user;
            }
        }
        return null;
    }
    
    /**
     * 创建新用户
     */
    public function create($data) {
        $id = count($this->users) + 1;
        
        $name = $data['name'] ?? '';
        $email = $data['email'] ?? '';
        $age = $data['age'] ?? 0;
        
        $user = new User($id, $name, $email, $age);
        $this->users[$id] = $user;
        
        return $user;
    }
    
    /**
     * 更新用户信息
     */
    public function update($id, $data) {
        if (!isset($this->users[$id])) {
            return null;
        }
        
        $user = $this->users[$id];
        
        if (isset($data['name'])) {
            $user->setName($data['name']);
        }
        if (isset($data['email'])) {
            $user->setEmail($data['email']);
        }
        if (isset($data['age'])) {
            $user->setAge($data['age']);
        }
        
        return $user;
    }
    
    /**
     * 删除用户
     */
    public function delete($id) {
        if (!isset($this->users[$id])) {
            return false;
        }
        
        unset($this->users[$id]);
        return true;
    }
    
    /**
     * 搜索用户
     */
    public function search($keyword, $field = 'name') {
        $results = [];
        
        foreach ($this->users as $user) {
            $value = '';
            switch ($field) {
                case 'name':
                    $value = $user->getName();
                    break;
                case 'email':
                    $value = $user->getEmail();
                    break;
            }
            
            if (stripos($value, $keyword) !== false) {
                $results[] = $user;
            }
        }
        
        return $results;
    }
}
