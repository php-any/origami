<?php

namespace tests\php;

/**
 * ComponentAttributeBag::class() 应对 CSS class 合并，不应变成 0="x" 属性。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$bag = new \Illuminate\View\ComponentAttributeBag(['class' => 'probe']);
$html = (string) $bag->class(['x']);

if (str_contains($html, '0=')) {
    Log::fatal('class() 误把列表当成属性: ' . $html);
}
if (!str_contains($html, 'class=') || !str_contains($html, 'x') || !str_contains($html, 'probe')) {
    Log::fatal('class() HTML 异常: ' . $html);
}

$assoc = (string) (new \Illuminate\View\ComponentAttributeBag())->class([
    'fi-sc',
    'fi-inline' => false,
    'fi-sc-has-gap' => true,
]);
if (str_contains($assoc, 'Object') || preg_match('/class="0[\s"]/', $assoc)) {
    Log::fatal('assoc class() 把数组 echo 成 Object/数字键: ' . $assoc);
}
if (!preg_match('/(?:^|[" ])fi-sc(?:[" ]|$)/', $assoc) || !str_contains($assoc, 'fi-sc-has-gap')) {
    Log::fatal('assoc class() 未合并 CSS class: ' . $assoc);
}

Log::info('ComponentAttributeBag::class 测试通过: ' . $html . ' | ' . $assoc);
