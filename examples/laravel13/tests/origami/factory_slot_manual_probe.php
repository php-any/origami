<?php
/**
 * Drive Factory startComponent/slot/endSlot/renderComponent manually.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$env = app('view');

// Use a simple anonymous view path
$view = 'filament::components.dropdown.index';

try {
    $env->startComponent($view, [
        'placement' => null,
        'teleport' => false,
        'attributes' => new Illuminate\View\ComponentAttributeBag([]),
    ]);
    echo "started\n";
    echo "ob_level=".ob_get_level()."\n";

    $env->slot('trigger', null, ['class' => 'my-trig']);
    echo "slot_started ob=".ob_get_level()."\n";
    echo '<button>OpenMe</button>';
    $env->endSlot();
    echo "endSlot ok\n";

    // Peek slots via reflection
    $ref = new ReflectionClass($env);
    // Factory uses trait; slots may be on class
    foreach (['slots','componentStack','slotStack'] as $prop) {
        try {
            $p = $ref->getProperty($prop);
            $p->setAccessible(true);
            $v = $p->getValue($env);
            echo "$prop=".json_encode(array_map(function ($x) {
                if (is_object($x)) return get_class($x);
                if (is_array($x)) return array_map(function ($y) {
                    return is_object($y) ? get_class($y) : $y;
                }, $x);
                return $x;
            }, is_array($v)?$v:[])). "\n";
        } catch (Throwable $e) {
            echo "$prop EX: ".$e->getMessage()."\n";
        }
    }

    // Also dump slots[0]['trigger'] detail
    $p = $ref->getProperty('slots');
    $p->setAccessible(true);
    $slots = $p->getValue($env);
    $t = $slots[0]['trigger'] ?? null;
    echo "trigger_type=".gettype($t).(is_object($t)?(' '.get_class($t)):'')."\n";
    if (is_object($t) && isset($t->attributes)) {
        echo "trigger_attrs_class=".get_class($t->attributes)."\n";
        echo "trigger_html=".$t->toHtml()."\n";
    }

    echo $env->renderComponent();
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}
echo "\nDONE\n";
