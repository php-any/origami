<?php

namespace tests\php;

/**
 * PHP preg_match：未匹配的可选捕获组不应出现在 $matches 中（ProgressBar isset($matches[2]) 依赖此语义）。
 */

$regex = '{%([a-z\-_]+)(?:\:([^%]+))?%}i';

$m = null;
preg_match($regex, '%current%', $m);
if (isset($m[2])) {
    Log::fatal('optional group unmatched must not isset($m[2]), got ' . var_export($m[2], true));
}
if (($m[1] ?? '') !== 'current') {
    Log::fatal('expected m[1]=current, got ' . var_export($m, true));
}

$m2 = null;
preg_match($regex, '%percent:3s%', $m2);
if (!isset($m2[2]) || $m2[2] !== '3s') {
    Log::fatal('expected m2[2]=3s, got ' . var_export($m2, true));
}

Log::info('preg_optional_group_test 测试通过');
