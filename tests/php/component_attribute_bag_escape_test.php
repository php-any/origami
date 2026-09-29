<?php

namespace tests\php;

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Illuminate\View\Compilers\BladeCompiler;
use Illuminate\View\ComponentAttributeBag;

$handler = "\$dispatch('close-modal', { id: 'database-notifications' })";
$sanitized = BladeCompiler::sanitizeComponentAttribute($handler);
$bag = new ComponentAttributeBag(['x-on:click' => $sanitized]);
$html = (string) $bag;

$iterator = $bag->getIterator();
if (! $iterator instanceof \Traversable) {
    Log::fatal('ComponentAttributeBag::getIterator() 应返回 Traversable，实际=' . get_debug_type($iterator));
}

if (str_contains($html, '&amp;#039;') || str_contains($html, '&amp;#39;')) {
    Log::fatal('ComponentAttributeBag 不应重复转义绑定属性: ' . $html);
}
if (! str_contains($html, '&#039;close-modal&#039;') && ! str_contains($html, '&#39;close-modal&#39;')) {
    Log::fatal('ComponentAttributeBag 丢失 Blade 已转义的属性值: ' . $html);
}

$booleans = (string) new ComponentAttributeBag([
    'disabled' => true,
    'x-data' => true,
    'wire:loading' => true,
    'hidden' => false,
]);
if ($booleans !== 'disabled="disabled" x-data="" wire:loading=""') {
    Log::fatal('ComponentAttributeBag 布尔属性错误: ' . $booleans);
}

Log::info('component_attribute_bag_escape 测试通过');
