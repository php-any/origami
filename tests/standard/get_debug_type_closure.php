<?php

namespace tests\standard;

/**
 * PHP 8：get_debug_type 对对象返回类名，对闭包返回 Closure。
 */

$fn = function () {
    return 1;
};
$t = get_debug_type($fn);
if ($t !== 'Closure' && strpos($t, 'Closure') === false) {
    Log::fatal('闭包 get_debug_type 失败: ' . $t);
}

Log::info('standard get_debug_type closure 测试通过');
