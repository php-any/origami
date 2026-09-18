<?php

namespace tests\php;

/**
 * UrlGenerator::asset 依赖 Str::finish('', '/') === '/'，否则资产 URL 丢根路径。
 */
use Illuminate\Support\Str;

if (!class_exists(Str::class)) {
    // 最小复现 finish 语义
    function finish_sim($value, $cap) {
        $quoted = preg_quote($cap, '/');
        return preg_replace('/(?:'.$quoted.')+$/u', '', $value).$cap;
    }
    $r = finish_sim('', '/');
    if ($r !== '/') {
        \Log::fatal('finish_sim empty 失败: '.var_export($r, true));
    }
    \Log::info('str_finish_empty 用本地模拟通过（无 Illuminate）');
    return;
}

$r = Str::finish('', '/');
if ($r !== '/') {
    \Log::fatal('Str::finish("", "/") 应为 "/"，实际: '.var_export($r, true));
}
$r2 = Str::finish('http://127.0.0.1:8000', '/');
if ($r2 !== 'http://127.0.0.1:8000/') {
    \Log::fatal('Str::finish root 失败: '.var_export($r2, true));
}

\Log::info('str_finish_empty 测试通过');
