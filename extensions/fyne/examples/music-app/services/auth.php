<?php
namespace MusicApp\Services;

/**
 * 认证服务 — 模拟登录
 */
class Auth {
    private $isLoggedIn = false;
    private $username = '';
    private $users = [
        'admin' => '123456',
        'user'  => 'password',
        'demo'  => 'demo',
    ];

    static function instance(): Auth {
        static $inst = null;
        if ($inst === null) {
            $inst = new Auth();
        }
        return $inst;
    }

    function tryLogin(string $username, string $password): bool {
        if (isset($this->users[$username]) && $this->users[$username] === $password) {
            $this->isLoggedIn = true;
            $this->username = $username;
            return true;
        }
        return false;
    }

    function logout(): void {
        $this->isLoggedIn = false;
        $this->username = '';
    }

    function isLoggedIn(): bool {
        return $this->isLoggedIn;
    }

    function getUsername(): string {
        return $this->username;
    }
}
