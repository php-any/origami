<?php

namespace tests\php;

/**
 * preg_replace 替换串中的 $1 / \\ 应与 PHP 一致（Symfony OutputFormatter::escape 依赖此行为）。
 */

$in = 'a<b>';
$out = preg_replace('/([^\\\\]|^)([<>])/', '$1\\\\$2', $in);
Log::info('escape result: ' . var_export($out, true));
Log::info('hex: ' . bin2hex($out));

// PHP: a\<b\>
$expected = 'a\\<b\\>';
if ($out !== $expected) {
    Log::fatal("expected " . var_export($expected, true) . " got " . var_export($out, true));
}

$out2 = preg_replace('/(foo)/', 'X$1Y', 'foo');
if ($out2 !== 'XfooY') {
    Log::fatal("simple \$1 failed: " . var_export($out2, true));
}

Log::info('preg_replace_backref_test 测试通过');
