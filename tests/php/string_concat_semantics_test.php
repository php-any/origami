<?php

namespace tests\php;

/**
 * 字符串拼接应对齐 PHP：int/bool/null 转字符串，且 . 可链式拼接。
 */
$s = 'n=' . 7 . true . null . 'x';
if ($s !== 'n=71x') {
    \Log::fatal('拼接结果错误: ' . var_export($s, true));
}

$a = 'ab';
$b = $a . $a;
if ($b !== 'abab') {
    \Log::fatal('字符串自拼接错误');
}

$lit = 'hello';
if ($lit !== 'hello') {
    \Log::fatal('字符串字面量读取错误');
}

\Log::info('string_concat_semantics 测试通过');
