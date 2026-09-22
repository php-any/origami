<?php

namespace tests\origami;

use Illuminate\View\ComponentAttributeBag;

if (!class_exists(ComponentAttributeBag::class, false)) {
    \Log::fatal('ComponentAttributeBag 应已预注册');
}

$bag = new ComponentAttributeBag(['id' => 'x', 'class' => 'a', 'wire:model' => 'name']);
if ($bag->get('id') !== 'x') {
    \Log::fatal('get 失败');
}
$html = (string) $bag;
if (strpos($html, 'id="x"') === false) {
    \Log::fatal('toHtml 失败: ' . $html);
}
$merged = $bag->merge(['role' => 'button']);
if ($merged->get('role') !== 'button') {
    \Log::fatal('merge 失败');
}
$only = $bag->whereStartsWith('wire:');
if (!$only->has('wire:model')) {
    \Log::fatal('whereStartsWith 失败');
}

\Log::info('illuminate_view_attribute_bag 测试通过');
