<?php

namespace tests\php;

/**
 * http_build_query 函数测试：将数组编码为 application/x-www-form-urlencoded 查询字符串。
 */

$basic = http_build_query(['a' => 1, 'b' => 2]);
if ($basic !== 'a=1&b=2') {
    Log::fatal('http_build_query 基本关联数组失败: ' . $basic);
}

$nested = http_build_query(['user' => ['name' => 'bob', 'age' => 10]]);
if ($nested !== 'user%5Bname%5D=bob&user%5Bage%5D=10') {
    Log::fatal('http_build_query 嵌套数组失败: ' . $nested);
}

$prefixed = http_build_query([0 => 'a', 1 => 'b'], 'item');
if ($prefixed !== 'item0=a&item1=b') {
    Log::fatal('http_build_query numeric_prefix 失败: ' . $prefixed);
}

$nullVal = http_build_query(['a' => null]);
if ($nullVal !== 'a=') {
    Log::fatal('http_build_query null 值应编码为空字符串: ' . $nullVal);
}

$space = http_build_query(['q' => 'hello world']);
if ($space !== 'q=hello+world') {
    Log::fatal('http_build_query RFC1738 空格应编码为 +: ' . $space);
}

$rfc3986 = http_build_query(['q' => 'hello world'], '', '&', 2);
if ($rfc3986 !== 'q=hello%20world') {
    Log::fatal('http_build_query RFC3986 编码失败: ' . $rfc3986);
}

Log::info('http_build_query 函数测试通过');
