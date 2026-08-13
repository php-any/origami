<?php

namespace tests\php;

/**
 * PHP 7.1+ 键名数组解构：['a' => $x] = $arr
 */
$route = [
	'action' => 'HomeController@index',
	'domain' => null,
	'method' => 'GET|HEAD',
	'middleware' => "web\napi",
	'uri' => '/',
];

[
	'action' => $action,
	'domain' => $domain,
	'method' => $method,
	'middleware' => $middleware,
	'uri' => $uri,
] = $route;

if ($action !== 'HomeController@index') {
	Log::fatal('keyed destructure action 失败: ' . var_export($action, true));
}
if ($method !== 'GET|HEAD') {
	Log::fatal('keyed destructure method 失败: ' . var_export($method, true));
}
if ($uri !== '/') {
	Log::fatal('keyed destructure uri 失败: ' . var_export($uri, true));
}
if ($middleware !== "web\napi") {
	Log::fatal('keyed destructure middleware 失败: ' . var_export($middleware, true));
}

Log::info('keyed array destructure 测试通过');
