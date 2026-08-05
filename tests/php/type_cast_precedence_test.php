<?php

/**
 * (int)/(string) 等类型转换的运算符优先级：高于字符串连接。
 */

$v = true;
if ((int) $v !== 1) {
    Log::fatal('(int)true 应为 1, got=' . (int) $v);
}

$s = "'" . (int) $v . "'";
if ($s !== "'1'") {
    Log::fatal("bool default 拼接失败: [$s]");
}

$s2 = (int) true . 'x';
if ($s2 !== '1x') {
    Log::fatal("(int)true.'x' 应为 1x, got=[$s2]");
}

Log::info('type cast 优先级测试通过');
