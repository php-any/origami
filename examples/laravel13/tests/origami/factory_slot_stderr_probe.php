<?php
/**
 * Factory slot lifecycle with STDERR diagnostics (avoid OB swallow).
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$env = app('view');
$log = function ($m) { fwrite(STDERR, $m."\n"); };

try {
    $env->startComponent('filament::components.dropdown.index', [
        'attributes' => new Illuminate\View\ComponentAttributeBag([]),
    ]);
    $log('started level='.ob_get_level());

    $env->slot('trigger', null, ['class' => 'my-trig']);
    $log('slot level='.ob_get_level());
    echo '<button>OpenMe</button>';
    $env->endSlot();
    $log('endSlot level='.ob_get_level());

    $ref = new ReflectionClass($env);
    $p = $ref->getProperty('slots');
    $p->setAccessible(true);
    $slots = $p->getValue($env);
    $log('slots_keys='.json_encode(array_keys($slots)));
    $log('slots0_keys='.json_encode(array_keys($slots[0] ?? [])));
    $t = $slots[0]['trigger'] ?? 'MISSING';
    $log('trigger='.(is_object($t) ? get_class($t) : var_export($t, true)));
    if (is_object($t)) {
        $log('html='.$t->toHtml());
        $log('attr='.get_class($t->attributes));
    }

    // Call componentData via reflection before render? render pops stack.
    // Instead render and catch
    $html = $env->renderComponent();
    $log('html_len='.strlen($html));
    $log('has_open='.(str_contains($html,'OpenMe')?'yes':'no'));
    echo $html;
} catch (Throwable $e) {
    $log('EX: '.$e->getMessage());
}
$log('DONE');
