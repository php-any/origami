<?php

namespace tests\php;

/**
 * dechex / hexdec 函数测试。
 */

$h = dechex(255);
if ($h !== 'ff') {
    Log::fatal('dechex(255) 失败: ' . $h);
}

$h0 = dechex(0);
if ($h0 !== '0') {
    Log::fatal('dechex(0) 失败: ' . $h0);
}

$n = hexdec('ff');
if ($n !== 255) {
    Log::fatal('hexdec(ff) 失败: ' . var_export($n, true));
}

$n2 = hexdec('0xff');
if ($n2 !== 255) {
    Log::fatal('hexdec(0xff) 失败: ' . var_export($n2, true));
}

$round = hexdec(dechex(12345));
if ($round !== 12345) {
    Log::fatal('dechex/hexdec 往返失败: ' . var_export($round, true));
}

// PHP：任意位置非法字符被忽略
$skip = hexdec('a-b');
if ($skip !== 171) {
    Log::fatal('hexdec(a-b) 失败: ' . var_export($skip, true));
}

Log::info('dechex/hexdec 函数测试通过');
