<?php

namespace tests\syntax;

/**
 * PHP 8.0：$obj::class 返回对象类名。
 */

class SyntaxObjClass_Demo
{
}

$o = new SyntaxObjClass_Demo();
$n = $o::class;
if ($n !== SyntaxObjClass_Demo::class && strpos($n, 'SyntaxObjClass_Demo') === false) {
    Log::fatal('$obj::class 失败: ' . var_export($n, true));
}

Log::info('syntax $obj::class 测试通过');
