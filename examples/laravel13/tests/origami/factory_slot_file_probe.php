<?php
/**
 * Factory slot lifecycle with file diagnostics.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$env = app('view');
$logFile = __DIR__.'/../../factory-slot-diag.txt';
@unlink($logFile);
$log = function ($m) use ($logFile) { file_put_contents($logFile, $m."\n", FILE_APPEND); };

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
    $log('slots_type='.gettype($slots));
    $log('slots='.var_export($slots, true));
    $t = is_array($slots) ? ($slots[0]['trigger'] ?? 'MISSING') : 'NOT_ARRAY';
    $log('trigger='.(is_object($t) ? (get_class($t).' html='.$t->toHtml()) : var_export($t, true)));

    $html = $env->renderComponent();
    $log('html_len='.strlen($html));
    $log('has_open='.(str_contains($html,'OpenMe')?'yes':'no'));
} catch (Throwable $e) {
    $log('EX: '.$e->getMessage().' @ '.$e->getFile().':'.$e->getLine());
}
$log('DONE');
echo file_get_contents($logFile);
