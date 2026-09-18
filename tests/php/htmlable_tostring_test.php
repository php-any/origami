<?php

namespace tests\php;

/**
 * (string) / e() 路径应对有 __toString 的 Htmlable 对象调用 __toString。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$bag = new \Illuminate\View\ComponentAttributeBag(['class' => 'x', 'id' => 'y']);

$cast = (string) $bag;
if (!str_contains($cast, 'class=') || str_contains($cast, 'Object(')) {
    Log::fatal('(string) ComponentAttributeBag 失败: ' . $cast);
}

if (!($bag instanceof \Illuminate\Contracts\Support\Htmlable)) {
    Log::fatal('instanceof Htmlable 失败');
}

$html = $bag->toHtml();
if (!str_contains($html, 'class=') || str_contains($html, 'Object(')) {
    Log::fatal('toHtml 失败: ' . $html);
}

$e = e($bag);
if (!str_contains($e, 'class=') || str_contains($e, 'Object(')) {
    Log::fatal('e() 失败: ' . $e);
}

Log::info('Htmlable/__toString 测试通过: ' . $e);
