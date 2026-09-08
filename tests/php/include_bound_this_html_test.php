<?php

namespace tests\php;

/**
 * Livewire 视图：Closure::bind 后 include 含 HTML 的模板，输出必须进入 ob 缓冲。
 */
class IncludeBoundThisHtml_Host
{
    public string $tag = 'from-this';
}

$host = new IncludeBoundThisHtml_Host();
$fn = \Closure::bind(function () {
    ob_start();
    include __DIR__ . '/include_bound_this_html_view.php';
    return ob_get_clean();
}, $host, $host);

$out = $fn();
if (!is_string($out) || !str_contains($out, '<div class="root">') || !str_contains($out, 'from-this')) {
    Log::fatal('bind+include HTML 未进入输出缓冲: ' . var_export($out, true));
}

Log::info('include HTML 绑定 $this 测试通过');
