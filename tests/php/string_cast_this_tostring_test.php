<?php

namespace tests\php;

/**
 * (string)$this 应调用 __toString，不能落到 Object(ClassName)。
 */

class ToString_ThisCastProbe
{
    public function __toString(): string
    {
        return 'from-toString';
    }

    public function asString(): string
    {
        return (string) $this;
    }
}

$o = new ToString_ThisCastProbe();
$direct = (string) $o;
if ($direct !== 'from-toString') {
    Log::fatal('(string)$obj 失败: ' . var_export($direct, true));
}
$viaThis = $o->asString();
if ($viaThis !== 'from-toString') {
    Log::fatal('(string)$this 失败: ' . var_export($viaThis, true));
}

Log::info('(string)$this __toString 测试通过');
