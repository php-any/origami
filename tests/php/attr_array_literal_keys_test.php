<?php

namespace tests\php;

/**
 * 数组字面量里嵌套静态调用 / 命名参数时必须保留 => 字符串键。
 * Blade withAttributes / icon-button merge 依赖此解析。
 */

class AttrLiteral_Sanitize
{
    public static function sanitize($v)
    {
        return $v;
    }
}

$withAttrs = [
    'color' => 'gray',
    'icon' => AttrLiteral_Sanitize::sanitize('bars'),
    'icon-alias' => AttrLiteral_Sanitize::sanitize('open'),
    'icon-size' => 'lg',
    'label' => AttrLiteral_Sanitize::sanitize('Expand sidebar'),
    'x-cloak' => true,
    'x-data' => '{}',
    'aria-controls' => 'fi-main-sidebar',
    'x-bind:aria-expanded' => '$store.sidebar.isOpen',
    'x-on:click' => '$store.sidebar.open()',
    'x-show' => '! $store.sidebar.isOpen',
    'class' => 'fi-topbar-open-sidebar-btn',
];

$keys = array_keys($withAttrs);
if ($keys !== [
    'color', 'icon', 'icon-alias', 'icon-size', 'label',
    'x-cloak', 'x-data', 'aria-controls', 'x-bind:aria-expanded',
    'x-on:click', 'x-show', 'class',
]) {
    Log::fatal('withAttributes 风格数组键丢失: ' . json_encode($keys));
}

$merge = [
    'aria-disabled' => null,
    // comment like compiled blade
    'aria-label' => AttrLiteral_Sanitize::sanitize('Expand sidebar', doubleEncode: false),
    'disabled' => false,
    'form' => null,
    'tabindex' => null,
    'type' => 'button',
    'wire:loading.attr' => 'disabled',
    'wire:target' => null,
];

$mkeys = array_keys($merge);
if ($mkeys !== [
    'aria-disabled', 'aria-label', 'disabled', 'form', 'tabindex',
    'type', 'wire:loading.attr', 'wire:target',
]) {
    Log::fatal('merge 风格数组键丢失: ' . json_encode($mkeys) . ' vals=' . json_encode(array_values($merge)));
}

function attr_literal_merge($arr, $escape = true)
{
    return array_keys($arr);
}

$named = attr_literal_merge([
    'aria-label' => 'Expand sidebar',
    'type' => 'button',
], escape: false);
if ($named !== ['aria-label', 'type']) {
    Log::fatal('命名参数 escape:false 后数组键丢失: ' . json_encode($named));
}

Log::info('属性数组字面量键测试通过');
