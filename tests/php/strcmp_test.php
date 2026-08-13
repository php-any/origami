<?php

namespace tests\php;

/**
 * strcmp / strcasecmp 内置函数注册冒烟。
 */
$r = strcmp('a', 'b');
if ($r >= 0) {
    Log::fatal('strcmp(a,b) 应 < 0，得 ' . var_export($r, true));
}
if (strcmp('b', 'a') <= 0) {
    Log::fatal('strcmp(b,a) 应 > 0');
}
if (strcmp('x', 'x') !== 0) {
    Log::fatal('strcmp(x,x) 应 === 0');
}

Log::info('strcmp 测试通过');
