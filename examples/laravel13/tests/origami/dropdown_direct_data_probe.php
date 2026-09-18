<?php
/**
 * Directly render dropdown compiled path with explicit trigger ComponentSlot in $data.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$trigger = new Illuminate\View\ComponentSlot('<button>OpenMe</button>', ['class' => 'my-trig']);
$slot = new Illuminate\View\ComponentSlot('PanelBody');

$data = [
    'trigger' => $trigger,
    'slot' => $slot,
    '__laravel_slots' => ['__default' => $slot, 'trigger' => $trigger],
    'attributes' => new Illuminate\View\ComponentAttributeBag([]),
    // defaults that @props may expect
    'availableHeight' => null,
    'availableWidth' => null,
    'flip' => true,
    'maxHeight' => null,
    'offset' => 8,
    'placement' => null,
    'shift' => false,
    'size' => false,
    'sizePadding' => 16,
    'teleport' => false,
    'width' => null,
];

try {
    $html = view('filament::components.dropdown.index', $data)->render();
    echo "len=".strlen($html)."\n";
    echo "has_open=".(str_contains($html,'OpenMe')?'yes':'no')."\n";
    echo "has_trig_class=".(str_contains($html,'fi-dropdown-trigger')?'yes':'no')."\n";
    echo "snip=".substr(preg_replace('/\s+/',' ',$html),0,300)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
