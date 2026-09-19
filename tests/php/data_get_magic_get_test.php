<?php

namespace tests\php;

/**
 * data_get($object, 'attr') 应对齐 PHP：走 __get，而不是只读声明属性。
 * Filament TextColumn 用 data_get($record, 'name') 取单元格。
 */

class DataGetMagic_User
{
    public function __get($name)
    {
        if ($name === 'email') {
            return 'admin@example.com';
        }
        if ($name === 'name') {
            return '管理员';
        }
        return null;
    }
}

if (!function_exists('data_get')) {
    \Log::info('data_get 未注册（非 Laravel helper 路径），跳过对象测试');
} else {
    $u = new DataGetMagic_User();
    $email = data_get($u, 'email');
    if ($email !== 'admin@example.com') {
        \Log::fatal('data_get 应走 __get, email='.var_export($email, true));
    }
    $name = data_get($u, 'name');
    if ($name !== '管理员') {
        \Log::fatal('data_get name 错误: '.var_export($name, true));
    }
}

\Log::info('data_get_magic_get 测试通过');
