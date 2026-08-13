<?php

namespace tests\php;

/**
 * 测试内置接口 Stringable 是否可用：
 *  - 类实现 \Stringable 并实现 __toString()
 *  - 类型提示 Stringable $x 是否能正常工作
 */

class Stringable_TestImpl implements \Stringable
{
    public function __toString(): string
    {
        return 'stringable-ok';
    }
}

function Stringable_accept(\Stringable $x): void
{
    Log::info('Stringable 接口测试: ' . (string) $x);
}

$obj = new Stringable_TestImpl();

if (!($obj instanceof \Stringable)) {
    Log::fatal('Stringable instanceof 测试失败');
}

Stringable_accept($obj);

Log::info('Stringable 接口测试通过');

