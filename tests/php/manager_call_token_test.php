<?php

namespace tests\php;

/**
 * 对齐 Laravel SessionManager::__call：csrf_token() 走 $this->driver()->$method(...$parameters)。
 * 空参数展开必须返回 driver 上 token() 的字符串，不能变成空串。
 */

class ManagerCall_TokenDriver
{
    public function token()
    {
        return 'ByoCftxDiqDOHtVQtTVlRpzCcmwt1vjpFZQAszdq';
    }
}

class ManagerCall_TokenManager
{
    public function driver()
    {
        return new ManagerCall_TokenDriver();
    }

    public function __call($method, $parameters)
    {
        return $this->driver()->$method(...$parameters);
    }
}

$manager = new ManagerCall_TokenManager();
$token = $manager->token();
if ($token !== 'ByoCftxDiqDOHtVQtTVlRpzCcmwt1vjpFZQAszdq') {
    Log::fatal('SessionManager 风格 __call+动态方法展开失败: ' . var_export($token, true));
}

$escaped = htmlspecialchars($token ?? '', ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8', true);
if ($escaped !== $token) {
    Log::fatal('e()/htmlspecialchars 把 CSRF token 变成了: ' . var_export($escaped, true));
}

Log::info('manager_call_token 测试通过');
