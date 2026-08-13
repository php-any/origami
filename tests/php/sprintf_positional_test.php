<?php

namespace tests\php;

/**
 * sprintf 位置参数 (%1$s / %2$s) 兼容性测试。
 *
 * 重点覆盖 TextDescriptor.php:71 的用法：
 *   sprintf($option->isNegatable() ? '--%1$s|--no-%1$s' : '--%1$s%2$s', $option->getName(), $value)
 */

// 情形 1：isNegatable = true => '--%1$s|--no-%1$s'
$name = 'help';
$value = '';
$fmtNeg = '--%1$s|--no-%1$s';
$synNeg = sprintf($fmtNeg, $name, $value);

if ($synNeg !== '--help|--no-help') {
    Log::fatal('sprintf 位置参数测试失败[neg]: 期望 --help|--no-help 实际 '.$synNeg);
}

// 情形 2：isNegatable = false => '--%1$s%2$s'
$fmtNon = '--%1$s%2$s';
$synNon = sprintf($fmtNon, $name, $value);

if ($synNon !== '--help') {
    Log::fatal('sprintf 位置参数测试失败[non-neg]: 期望 --help 实际 '.$synNon);
}

Log::info('sprintf 位置参数测试通过: '.$synNeg.' / '.$synNon);

