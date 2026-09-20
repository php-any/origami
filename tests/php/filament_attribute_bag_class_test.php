<?php

namespace tests\php;

/**
 * Filament ComponentAttributeBag::class 必须把新 class 并入已有 class 键，
 * 不得留下 0="fi-body-has-navigation"。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

$bag = new \Filament\Support\View\ComponentAttributeBag([
    'class' => 'fi-body-has-navigation fi-body-has-topbar',
    'livewire' => null,
]);

$html = (string) $bag->merge([], escape: false)->class([
    'fi-body',
    'fi-panel-admin',
]);

if (str_contains($html, ' 0="') || preg_match('/(?:^|\s)0="/', $html)) {
    Log::fatal('Filament bag class 数字键: ' . $html);
}
if (!str_contains($html, 'class="') || !str_contains($html, 'fi-body-has-navigation') || !str_contains($html, 'fi-body') || !str_contains($html, 'fi-panel-admin')) {
    Log::fatal('Filament bag class 未合并: ' . $html);
}

$illuminate = new \Illuminate\View\ComponentAttributeBag([
    'class' => 'fi-body-has-navigation fi-body-has-topbar',
]);
$html2 = (string) $illuminate->merge([], escape: false)->class([
    'fi-body',
    'fi-panel-admin',
]);
if (str_contains($html2, ' 0="') || preg_match('/(?:^|\s)0="/', $html2)) {
    Log::fatal('Illuminate bag class 数字键: ' . $html2);
}
if (!str_contains($html2, 'fi-body-has-navigation') || !str_contains($html2, 'fi-body')) {
    Log::fatal('Illuminate bag class 未合并: ' . $html2);
}

Log::info('Filament/Illuminate bag class 测试通过: ' . $html . ' || ' . $html2);
