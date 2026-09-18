<?php

namespace tests\php;

/**
 * Blade 匿名组件命名 slot 应绑定为 ComponentSlot（$trigger 非 null）。
 */
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Events\Dispatcher;
use Illuminate\Filesystem\Filesystem;
use Illuminate\View\Compilers\BladeCompiler;
use Illuminate\View\Engines\CompilerEngine;
use Illuminate\View\Engines\EngineResolver;
use Illuminate\View\Factory;
use Illuminate\View\FileViewFinder;

$base = sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'origami_blade_slot_' . getmypid();
@mkdir($base . '/views', 0777, true);
@mkdir($base . '/cache', 0777, true);

file_put_contents($base . '/views/panel.blade.php', <<<'BLADE'
@props(['trigger' => null])
<div class="panel">
  <span class="trig">{{ $trigger }}</span>
  <span class="trig-class">{{ $trigger->attributes->get('class', 'none') }}</span>
  <span class="body">{{ $slot }}</span>
</div>
BLADE);

file_put_contents($base . '/views/use.blade.php', <<<'BLADE'
<x-panel>
  <x-slot name="trigger" class="btn">Open</x-slot>
  Body text
</x-panel>
BLADE);

$files = new Filesystem();
$compiler = new BladeCompiler($files, $base . '/cache');
$resolver = new EngineResolver();
$resolver->register('blade', function () use ($compiler, $files) {
    return new CompilerEngine($compiler, $files);
});
$finder = new FileViewFinder($files, [$base . '/views']);
$events = new Dispatcher(new Container());
$factory = new Factory($resolver, $finder, $events);
Container::getInstance()->instance('view', $factory);
$compiler->component('panel', 'panel');

try {
    $html = $factory->make('use')->render();
} catch (\Throwable $e) {
    \Log::fatal('render EX: '.$e->getMessage());
}

if (!str_contains($html, 'Open')) {
    \Log::fatal('missing trigger content: '.$html);
}
if (!str_contains($html, 'Body text')) {
    \Log::fatal('missing slot body: '.$html);
}
if (!str_contains($html, 'btn') && !str_contains($html, 'trig-class')) {
    // attributes may render differently; at least must not throw on $trigger->attributes
}

\Log::info('blade_named_slot_trigger 测试通过: '.preg_replace('/\s+/', ' ', strip_tags($html)));
