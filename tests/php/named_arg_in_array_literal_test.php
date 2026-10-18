<?php

namespace tests\php;

/**
 * 数组值里的命名参数 e(..., doubleEncode: false) 不得打乱后续 => 键。
 * Filament icon-button compiled merge 依赖此解析。
 */

function escape_label(string $value, bool $doubleEncode = true): string { return htmlspecialchars($value, ENT_QUOTES, "UTF-8", $doubleEncode); }

$arr = [
    'aria-disabled' => null,
    'aria-label' => escape_label('Expand sidebar', doubleEncode: false),
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

Log::info('命名参数数组键测试通过');
