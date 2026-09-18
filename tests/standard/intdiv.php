<?php

namespace tests\standard;

/**
 * PHP 7+ / 8：intdiv 整数除法。
 */

if (intdiv(10, 3) !== 3) {
    Log::fatal('intdiv(10,3) 失败');
}
if (intdiv(-10, 3) !== -3) {
    Log::fatal('intdiv 向零取整失败: ' . intdiv(-10, 3));
}

Log::info('standard intdiv 测试通过');
