<?php

namespace tests\php;

/**
 * 静态方法继承调用测试：
 * - 父类 Base 定义 public static function now()
 * - 子类 Child extends Base
 * - Child::now() 应该能正常调用（覆盖之前静态方法不查父类的问题）
 */

class StaticBase
{
    public static function now(): string
    {
        return 'base-now';
    }
}

class StaticChild extends StaticBase
{
}

$val = StaticChild::now();

if ($val !== 'base-now') {
    Log::fatal('StaticChild::now() 调用失败，返回值: ' . var_export($val, true));
}

Log::info('静态方法继承调用测试通过');

