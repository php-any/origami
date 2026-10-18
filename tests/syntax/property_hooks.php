<?php

namespace tests\syntax;

/**
 * PHP 8.4：property hooks 必须执行，且每次赋值只调用一次。
 */

class SyntaxHooks_C
{
    public int $writes = 0;
    public string $name {
        set {
            $this->writes++;
            $this->name = strtoupper($value);
        }
    }
}

$o = new SyntaxHooks_C();
$o->name = 'x';
if ($o->name !== 'X' || $o->writes !== 1) {
    throw new \Exception('property set hook must run exactly once');
}
$o->name = 'second';
if ($o->name !== 'SECOND' || $o->writes !== 2) {
    throw new \Exception('property set hook must run on subsequent writes');
}

echo "property hook writes OK\n";
