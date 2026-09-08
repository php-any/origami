<?php

namespace tests\php;

/**
 * Closure::bind 后 include 必须能使用绑定对象的 $this（Livewire Blade 引擎依赖此语义）。
 */
class IncludeBoundThis_Host
{
    public string $tag = 'from-this';
}

$host = new IncludeBoundThis_Host();
$fn = \Closure::bind(function () {
    ob_start();
    include __DIR__ . '/include_bound_this_view.php';
    return ob_get_clean();
}, $host, $host);

$out = $fn();
if ($out !== 'from-this') {
    Log::fatal('bind+include 未能读到 $this: ' . var_export($out, true));
}

Log::info('include 绑定 $this 测试通过');
