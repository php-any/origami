<?php

namespace tests\php;

/**
 * PHP 8.4 property hooks：写入必须执行 set 钩子。
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
if ($d->name !== 'HELLO') {
    throw new \Exception('set hook must uppercase the stored value');
}

echo "property hook set OK\n";
