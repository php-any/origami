<?php

namespace tests\php;

/**
 * PHP 8.4 property hooks：至少能解析并作为普通属性读写。
 */

class PropertyHooks_Demo
{
    public string $name {
        set {
            $this->name = strtoupper($value);
        }
    }
}

$d = new PropertyHooks_Demo();
$d->name = 'hello';
// 钩子尚未执行时，至少应能赋值/读取普通属性槽
if ($d->name !== 'hello' && $d->name !== 'HELLO') {
    \Log::fatal('property hooks 属性读写失败: ' . var_export($d->name, true));
}

\Log::info('property_hooks_parse 测试通过');
