<?php

namespace tests\syntax;

/**
 * PHP 8.4：property hooks 至少可解析并读写。
 */

class SyntaxHooks_C
{
    public string $name {
        set {
            $this->name = $value;
        }
    }
}

$o = new SyntaxHooks_C();
$o->name = 'x';
if ($o->name !== 'x' && $o->name !== 'X') {
    Log::fatal('property hooks 读写失败');
}

Log::info('syntax property hooks 测试通过');
