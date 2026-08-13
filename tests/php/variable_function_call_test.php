<?php

namespace tests\php;

/**
 * 验证可变函数调用 $fn()（字符串函数名）。
 */
function __var_fn_probe() {
	return 'ok';
}

$handler = 'tests\\php\\__var_fn_probe';
$ret = $handler();
if ($ret !== 'ok') {
	Log::fatal('字符串可变函数调用失败: ' . var_export($ret, true));
}

$routes = [
	'login' => 'tests\\php\\__var_fn_probe',
];
$route = 'login';
$h = $routes[$route] ?? null;
if ($h === null) {
	Log::fatal('关联数组取路由 handler 失败');
}
$ret2 = $h();
if ($ret2 !== 'ok') {
	Log::fatal('从路由表取出的 handler 调用失败: ' . var_export($ret2, true));
}

Log::info('variable function call 测试通过');
