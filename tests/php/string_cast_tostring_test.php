<?php

namespace tests\php;

/**
 * (string) 强转应调用对象的 __toString。
 */

class StringCast_ToStringable
{
    public function __toString(): string
    {
        return 'cast-ok';
    }
}

$o = new StringCast_ToStringable();
$s = (string) $o;
if ($s !== 'cast-ok') {
    Log::fatal('(string) cast __toString failed: ' . var_export($s, true));
}

Log::info('(string) cast __toString 测试通过');
