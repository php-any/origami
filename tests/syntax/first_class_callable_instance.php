<?php

namespace tests\syntax;

/**
 * PHP 8.1：实例方法一等可调用 $obj->m(...)。
 */

class SyntaxFcc_Host
{
    public function twice(int $n): int
    {
        return $n * 2;
    }
}

$h = new SyntaxFcc_Host();
$fn = $h->twice(...);
if ($fn(4) !== 8) {
    Log::fatal('实例 first-class callable 失败');
}

Log::info('syntax instance first-class callable 测试通过');
