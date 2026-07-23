<?php

namespace tests\php;

/**
 * preg_replace 支持 \\1 / $1 反引用；preg_match 未匹配的尾部可选组不可 isset。
 */

// 1) OutputWrapper 风格：'\\1' + "\n"
$out = preg_replace('/(foo)/', '\\1!', 'foo');
if ($out !== 'foo!') {
    Log::fatal('backslash-1 backref failed: ' . var_export($out, true));
}

$out2 = preg_replace('/(a)(b)/', '$2-$1', 'ab');
if ($out2 !== 'b-a') {
    Log::fatal('$n backref failed: ' . var_export($out2, true));
}

// OutputFormatter::escape 风格
$esc = preg_replace('/([^\\\\]|^)([<>])/', '$1\\\\$2', 'a<b>');
if ($esc !== 'a\\<b\\>') {
    Log::fatal('escape-style replace failed: ' . var_export($esc, true) . ' hex=' . bin2hex($esc));
}

// 2) 可选捕获组
$regex = '{%([a-z\-_]+)(?:\:([^%]+))?%}i';
$m = null;
preg_match($regex, '%current%', $m);
if (isset($m[2])) {
    Log::fatal('unmatched optional group must not isset, got ' . var_export($m, true));
}
preg_match($regex, '%percent:3s%', $m2);
if (!isset($m2[2]) || $m2[2] !== '3s') {
    Log::fatal('matched optional group missing: ' . var_export($m2, true));
}

Log::info('preg_php_backref_optional_test 测试通过');
