<?php

namespace tests\php;

/**
 * hash() 第二参数为 null 时按 PHP 强制为 ""（空串摘要）。
 */

$a = hash('sha256', '');
$b = hash('sha256', null);
if ($a !== $b) {
    Log::fatal('hash(null) 应与 hash("") 相同');
}
if (strlen($b) !== 64) {
    Log::fatal('hash sha256 长度应为 64');
}

Log::info('hash null 强制转换测试通过');
