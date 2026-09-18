<?php
/**
 * Minimal Blade named-slot trigger under laravel13.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$dir = storage_path('framework/views');
$cache = storage_path('framework/cache');
@mkdir($dir, 0777, true);

$viewPath = resource_path('views');
@mkdir($viewPath.'/origami_slot_test', 0777, true);
file_put_contents($viewPath.'/origami_slot_test/panel.blade.php', <<<'BLADE'
@props(['trigger' => null])
<div>
@if($trigger)
TRIG={{ $trigger }}|
ATTR={{ $trigger->attributes->get('class') }}|
@else
TRIG_NULL
@endif
SLOT={{ $slot }}
</div>
BLADE);

file_put_contents($viewPath.'/origami_slot_test/use.blade.php', <<<'BLADE'
<x-origami_slot_test.panel>
  <x-slot name="trigger" class="btn">Open</x-slot>
  Body
</x-origami_slot_test.panel>
BLADE);

try {
    $html = view('origami_slot_test.use')->render();
    echo "html=".preg_replace('/\s+/',' ', $html)."\n";
    echo "has_open=".(str_contains($html,'Open')?'yes':'no')."\n";
    echo "trig_null=".(str_contains($html,'TRIG_NULL')?'yes':'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
