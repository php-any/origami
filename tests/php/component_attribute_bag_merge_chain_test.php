<?php

namespace tests\php;

/**
 * Filament icon-button：merge(escape: false)+class() 后 HTML 必须是属性名而非 0=。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Illuminate\View\ComponentAttributeBag;

$propDefaults = [
    'badge' => null,
    'badgeColor' => 'primary',
    'badgeSize' => null,
    'color' => 'primary',
    'disabled' => false,
    'form' => null,
    'formId' => null,
    'href' => null,
    'icon' => null,
    'iconAlias' => null,
    'iconSize' => null,
    'keyBindings' => null,
    'label' => null,
    'loadingIndicator' => true,
    'size' => null,
    'spaMode' => null,
    'tag' => 'button',
    'target' => null,
    'tooltip' => null,
    'type' => 'button',
];

$incoming = [
    'color' => 'gray',
    'icon' => 'bars',
    'icon-alias' => 'open',
    'icon-size' => 'lg',
    'label' => 'Expand sidebar',
    'x-cloak' => true,
    'x-data' => '{}',
    'aria-controls' => 'fi-main-sidebar',
    'x-bind:aria-expanded' => '$store.sidebar.isOpen',
    'x-on:click' => '$store.sidebar.open()',
    'x-show' => '! $store.sidebar.isOpen',
    'class' => 'fi-topbar-open-sidebar-btn',
];

$__propNames = ComponentAttributeBag::extractPropNames($propDefaults);
$__newAttributes = [];
foreach ($incoming as $__key => $__value) {
    if (in_array($__key, $__propNames)) {
        continue;
    }
    $__newAttributes[$__key] = $__value;
}

$keys = array_keys($__newAttributes);
if (!in_array('class', $keys, true) || !in_array('x-data', $keys, true) || in_array(0, $keys, true)) {
    Log::fatal('extract 后属性键丢失: ' . json_encode($keys));
}

$attributes = new ComponentAttributeBag($__newAttributes);
$html = (string) $attributes
    ->merge([
        'aria-disabled' => null,
        'aria-label' => 'Expand sidebar',
        'disabled' => false,
        'form' => null,
        'tabindex' => null,
        'type' => 'button',
        'wire:loading.attr' => 'disabled',
        'wire:target' => null,
    ], escape: false)
    ->merge([
        'title' => 'Expand sidebar',
    ], escape: true)
    ->class([
        'fi-icon-btn',
        'fi-disabled' => false,
        'fi-size-md',
    ]);

if (str_contains($html, ' 0="') || preg_match('/\s0="/', $html)) {
    Log::fatal('icon-button merge 链把属性变成数字键: ' . $html);
}
if (!str_contains($html, 'class="') || !str_contains($html, 'fi-icon-btn') || !str_contains($html, 'fi-topbar-open-sidebar-btn')) {
    Log::fatal('icon-button class 未合并: ' . $html);
}
if (!str_contains($html, 'aria-label="Expand sidebar"') || !str_contains($html, 'type="button"') || !str_contains($html, 'x-data="{}"')) {
    Log::fatal('icon-button 具名属性丢失: ' . $html);
}

// Filament color(gray) 对 IconButton 会 class([]) → merge(['class' => ''])，不得把 class 冲成 0=
$afterEmptyClass = (string) $attributes
    ->merge([
        'aria-label' => 'Expand sidebar',
        'type' => 'button',
    ], escape: false)
    ->class([
        'fi-icon-btn',
        'fi-size-md',
        'fi-topbar-open-sidebar-btn',
    ])
    ->class([]);
if (str_contains($afterEmptyClass, ' 0="') || preg_match('/\s0="/', $afterEmptyClass)) {
    Log::fatal('class([]) 把属性变成数字键: ' . $afterEmptyClass);
}
if (!str_contains($afterEmptyClass, 'fi-icon-btn')) {
    Log::fatal('class([]) 冲掉了 class: ' . $afterEmptyClass);
}

Log::info('ComponentAttributeBag merge 链测试通过: ' . $html);
