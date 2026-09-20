<?php

namespace tests\php;

/**
 * layout.base body：withAttributes 的 class 来自 Arr::toCssClasses + sanitize，
 * 再 class(['fi-body', ...]) 不得变成 class="fi-body" 0="fi-body-has-navigation"。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$hasNavigation = true;
$isSidebarCollapsibleOnDesktop = false;
$isSidebarFullyCollapsibleOnDesktop = false;
$hasTopbar = true;
$hasTopNavigation = false;
$livewire = null;

$incoming = [
    'livewire' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($livewire),
    'class' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Illuminate\Support\Arr::toCssClasses([
        'fi-body-has-navigation' => $hasNavigation,
        'fi-body-has-sidebar-collapsible-on-desktop' => $isSidebarCollapsibleOnDesktop,
        'fi-body-has-sidebar-fully-collapsible-on-desktop' => $isSidebarFullyCollapsibleOnDesktop,
        'fi-body-has-topbar' => $hasTopbar,
        'fi-body-has-top-navigation' => $hasTopNavigation,
    ])),
];

$keys = array_keys($incoming);
if ($keys !== ['livewire', 'class']) {
    Log::fatal('withAttributes 字面量键丢失: ' . json_encode($keys) . ' vals=' . json_encode($incoming));
}

$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(['livewire' => null]);
$__newAttributes = [];
foreach ($incoming as $__key => $__value) {
    if (in_array($__key, $__propNames)) {
        continue;
    }
    $__newAttributes[$__key] = $__value;
}
$nk = array_keys($__newAttributes);
if ($nk !== ['class']) {
    Log::fatal('extract 后 class 键丢失: keys=' . json_encode($nk) . ' props=' . json_encode($__propNames));
}

$attributes = new \Illuminate\View\ComponentAttributeBag($__newAttributes);
$html = (string) $attributes
    ->merge([], escape: false)
    ->class([
        'fi-body',
        'fi-panel-admin',
    ]);

if (str_contains($html, ' 0="') || preg_match('/\s0="/', $html)) {
    Log::fatal('body class 链数字键: ' . $html);
}
if (!str_contains($html, 'fi-body') || !str_contains($html, 'fi-panel-admin') || !str_contains($html, 'fi-body-has-navigation') || !str_contains($html, 'fi-body-has-topbar')) {
    Log::fatal('body class 未合并: ' . $html);
}

Log::info('layout.base body 属性测试通过: ' . $html);
