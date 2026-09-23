<?php

namespace tests\php;

/**
 * ComponentAttributeBag：extractPropNames / when / merge 热路径冒烟。
 */

if (!class_exists(\Illuminate\View\ComponentAttributeBag::class, false)) {
    Log::info('skip: ComponentAttributeBag 原生类未注册');
    return;
}

$props = ['wire:model', 'fooBar' => true];
$names = \Illuminate\View\ComponentAttributeBag::extractPropNames($props);
if (!is_array($names)) {
    Log::fatal('extractPropNames 返回类型错误: ' . gettype($names));
}
$flat = array_values($names);
if (!in_array('wire:model', $flat, true)) {
    Log::fatal('extractPropNames 缺 wire:model: ' . json_encode($flat));
}
if (!in_array('fooBar', $flat, true) || !in_array('foo-bar', $flat, true)) {
    Log::fatal('extractPropNames 未生成 kebab 变体: ' . json_encode($flat));
}

$bag = new \Illuminate\View\ComponentAttributeBag(['class' => 'a', 'id' => 'x']);
$merged = $bag->merge(['data-x' => '1']);
$html = (string) $merged;
if (!str_contains($html, 'class="a"') || !str_contains($html, 'data-x="1"')) {
    Log::fatal('merge/toString 异常: ' . $html);
}

$when = $bag->when(true, function ($b) {
    return $b->merge(['ok' => '1']);
});
if (!str_contains((string) $when, 'ok="1"')) {
    Log::fatal('when(true) 未应用 callback');
}

Log::info('illuminate_view_attribute_bag 测试通过');
