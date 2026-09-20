<?php

namespace tests\php;

/**
 * 数组值里的命名参数 e(..., doubleEncode: false) 不得打乱后续 => 键。
 * Filament icon-button compiled merge 依赖此解析。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$arr = [
    'aria-disabled' => null,
    'aria-label' => e('Expand sidebar', doubleEncode: false),
    'disabled' => false,
    'form' => null,
    'tabindex' => null,
    'type' => 'button',
    'wire:loading.attr' => 'disabled',
    'wire:target' => null,
];

$keys = array_keys($arr);
if ($keys !== [
    'aria-disabled', 'aria-label', 'disabled', 'form', 'tabindex',
    'type', 'wire:loading.attr', 'wire:target',
]) {
    Log::fatal('e(doubleEncode:) 扰乱数组键: ' . json_encode($keys) . ' vals=' . json_encode($arr));
}

$bag = new \Illuminate\View\ComponentAttributeBag([
    'x-cloak' => true,
    'x-data' => '{}',
    'class' => 'fi-topbar-open-sidebar-btn',
]);
$html = (string) $bag->merge($arr, escape: false)->class(['fi-icon-btn', 'fi-size-md']);
if (str_contains($html, ' 0="') || preg_match('/\s0="/', $html)) {
    Log::fatal('带 doubleEncode 的 merge 变成数字键: ' . $html);
}
if (!str_contains($html, 'aria-label="Expand sidebar"') || !str_contains($html, 'class="')) {
    Log::fatal('带 doubleEncode 的 merge HTML 异常: ' . $html);
}

Log::info('e(doubleEncode) 数组键测试通过: ' . $html);
