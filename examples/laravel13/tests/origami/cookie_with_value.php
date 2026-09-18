<?php

// Cookie::withValue 对齐 Symfony：EncryptCookies::duplicate 依赖此方法。

use Symfony\Component\HttpFoundation\Cookie;

require __DIR__.'/../../vendor/autoload.php';

$c = new Cookie('laravel_session', 'plain', 0, '/', null, false, true);
$n = $c->withValue('encrypted-payload');
if ($n->getName() !== 'laravel_session') {
    echo "FAIL: withValue 应保留 name\n";
    exit(1);
}
if ($n->getValue() !== 'encrypted-payload') {
    echo "FAIL: withValue 应替换 value, 实际=".$n->getValue()."\n";
    exit(1);
}
if ($c->getValue() !== 'plain') {
    echo "FAIL: withValue 不应改原对象\n";
    exit(1);
}

echo "OK: cookie_with_value\n";
