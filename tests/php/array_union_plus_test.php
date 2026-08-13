<?php

namespace tests\php;

/**
 * 验证 PHP 数组 + 运算符按键并集且保留字符串键（route_url 依赖）。
 */
$a = 'login';
$params = [];
$x = ['a' => $a] + $params;
if (!array_key_exists('a', $x) || $x['a'] !== 'login') {
	Log::fatal('["a"=>$a]+[] 应保留字符串键 a，实际: ' . var_export($x, true));
}
$q = http_build_query($x, '', '&', PHP_QUERY_RFC3986);
if ($q !== 'a=login') {
	Log::fatal('http_build_query 期望 a=login，实际: ' . $q);
}

$y = ['a' => 'login', 0 => 'keep'] + [0 => 'drop', 'b' => 'add', 1 => 'one'];
if (($y['a'] ?? null) !== 'login' || ($y[0] ?? null) !== 'keep' || ($y['b'] ?? null) !== 'add' || ($y[1] ?? null) !== 'one') {
	Log::fatal('数组并集键优先语义错误: ' . var_export($y, true));
}

Log::info('array + union 测试通过');
