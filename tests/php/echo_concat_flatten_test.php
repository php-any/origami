<?php

namespace tests\php;

/**
 * echo 字面量、变量、`.` 拼接与 __toString 必须与原先输出一致。
 */

class EchoFlatten_ToString
{
    public function __toString(): string
    {
        return 'OBJ';
    }
}

ob_start();
$x = 'mid';
echo 'A'.$x.'B';
echo 7;
echo true;
echo new EchoFlatten_ToString();
$got = ob_get_clean();

if ($got !== 'AmidB71OBJ') {
    \Log::fatal('echo 拼接输出不符: '.var_export($got, true));
}

$y = 'L'.'R';
if ($y !== 'LR') {
    \Log::fatal('赋值拼接不符: '.var_export($y, true));
}

\Log::info('echo concat flatten 测试通过');
