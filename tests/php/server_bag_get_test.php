<?php

namespace tests\php;

/**
 * Symfony ServerBag 继承 ParameterBag::get/has（Telescope 用 $request->server->get）。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$bag = new \Symfony\Component\HttpFoundation\ServerBag([
    'QUERY_STRING' => 'a=1',
    'REQUEST_URI' => '/telescope',
]);

$qs = $bag->get('QUERY_STRING');
if ($qs !== 'a=1') {
    Log::fatal('ServerBag::get QUERY_STRING 失败: ' . var_export($qs, true));
}
if (!$bag->has('REQUEST_URI')) {
    Log::fatal('ServerBag::has REQUEST_URI 应为 true');
}
$missing = $bag->get('MISSING_KEY', 'def');
if ($missing !== 'def') {
    Log::fatal('ServerBag::get default 失败: ' . var_export($missing, true));
}

$input = new \Symfony\Component\HttpFoundation\InputBag(['q' => 'ok']);
if ($input->get('q') !== 'ok') {
    Log::fatal('InputBag::get 失败');
}

Log::info('ServerBag 继承 ParameterBag::get 测试通过');
