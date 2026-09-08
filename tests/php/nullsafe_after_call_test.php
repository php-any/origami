<?php

namespace tests\php;

/**
 * 函数/方法调用后的 ?->（后台布局 auth()->user()?->hasPermission）。
 */
class NullsafeAfterCall_User
{
    public $name = 'admin';

    public function hasPermission($p)
    {
        return $p === 'admins.view';
    }
}

class NullsafeAfterCall_Guard
{
    public function user()
    {
        return new NullsafeAfterCall_User();
    }
}

function nullsafe_after_call_auth($g)
{
    return new NullsafeAfterCall_Guard();
}

if (!nullsafe_after_call_auth('admin')->user()?->hasPermission('admins.view')) {
    Log::fatal('auth()->user()?->hasPermission 应为 true');
}

$n = nullsafe_after_call_auth('admin')->user()?->name;
if ($n !== 'admin') {
    Log::fatal('auth()->user()?->name 失败: ' . var_export($n, true));
}

Log::info('调用后 ?-> 测试通过');
