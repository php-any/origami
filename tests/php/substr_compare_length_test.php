<?php

namespace tests\php;

/**
 * PHP substr_compare 在指定 length 时只比较双方前 length 个字符。
 * Symfony IpUtils::checkIp4 用它做 CIDR 前缀匹配。
 */
$full = substr_compare('abcdef', 'cde', 2);
if ($full !== 1) {
	Log::fatal('substr_compare 省略 length 应按 PHP 比较剩余 haystack, 实际: ' . var_export($full, true));
}
$prefix = substr_compare('abcdefgh', 'abcXXXXX', 0, 3);
if ($prefix !== 0) {
	Log::fatal('指定 length=3 应只比较前缀, 实际: ' . var_export($prefix, true));
}
$mismatch = substr_compare('abcdefgh', 'abdXXXXX', 0, 3);
if ($mismatch === 0) {
	Log::fatal('前缀不同时应非 0');
}
Log::info('substr_compare length 测试通过');
