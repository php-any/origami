<?php

namespace tests\php;

/**
 * trait 静态方法测试：
 * - TraitStatic 提供 public static function fromTrait()
 * - ClassStatic 使用 use TraitStatic;
 * - ClassStatic::fromTrait() 应该可以正常调用。
 */

trait TraitStatic
{
    public static function fromTrait(): string
    {
        return 'trait-static';
    }
}

class ClassStatic
{
    use TraitStatic;
}

$v = ClassStatic::fromTrait();
if ($v !== 'trait-static') {
    Log::fatal('ClassStatic::fromTrait() 调用失败，返回值: ' . var_export($v, true));
}

Log::info('trait 静态方法测试通过');

