<?php

namespace tests\php;

/**
 * Laravel RouteUrlGenerator::replaceNamedParameters 依赖
 * preg_replace_callback + use (&$parameters) 替换 {orderId}。
 */
$path = 'admin/orders/{orderId}';
$parameters = ['orderId' => 1];
$n = preg_match_all('/\{(.*?)(\?)?\}/', $path, $m);
if ($n !== 1) {
    Log::fatal('preg_match_all 未匹配路由占位符 n=' . var_export($n, true) . ' m=' . var_export($m, true));
}
if (($m[1][0] ?? '') !== 'orderId') {
    Log::fatal('捕获组 orderId 错误: ' . var_export($m, true));
}

$replaced = preg_replace_callback('/\{(.*?)(\?)?\}/', function ($match) use (&$parameters) {
    $name = $match[1];
    if (isset($parameters[$name]) && $parameters[$name] !== '') {
        $v = $parameters[$name];
        unset($parameters[$name]);
        return $v;
    }
    return $match[0];
}, $path);

if ($replaced !== 'admin/orders/1') {
    Log::fatal('preg_replace_callback 未替换路由参数: ' . var_export($replaced, true) . ' leftover=' . var_export($parameters, true));
}

Log::info('preg_replace_callback 路由参数替换测试通过');
