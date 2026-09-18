<?php
/**
 * Dump whether trigger is in data before/after @props simulation for dropdown.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$env = app('view');
$logFile = __DIR__.'/../../factory-slot-diag2.txt';
@unlink($logFile);
$log = function ($m) use ($logFile) { file_put_contents($logFile, $m."\n", FILE_APPEND); };

$env->startComponent('filament::components.dropdown.index', [
    'attributes' => new Illuminate\View\ComponentAttributeBag([]),
]);
$env->slot('trigger', null, ['class' => 'my-trig']);
echo '<button>OpenMe</button>';
$env->endSlot();

// Reflect and call componentData-like merge by rendering a spy view instead.
// Hook: temporarily replace evaluate by dumping data via View composer? 
// Simpler: reflection invoke protected componentData after popping manually.

$ref = new ReflectionClass($env);
$stack = $ref->getProperty('componentStack');
$stack->setAccessible(true);
$dataProp = $ref->getProperty('componentData');
$dataProp->setAccessible(true);
$slotsProp = $ref->getProperty('slots');
$slotsProp->setAccessible(true);

$log('stack_count_before_pop='.count($stack->getValue($env)));
$view = array_pop($stack->getValue($env));
// array_pop on getValue may not write back — use setValue
$s = $stack->getValue($env);
array_pop($s);
$stack->setValue($env, $s);
$log('stack_count_after_pop='.count($stack->getValue($env)));

$method = $ref->getMethod('componentData');
$method->setAccessible(true);
try {
    $data = $method->invoke($env);
    $log('data_keys='.json_encode(array_keys($data)));
    $log('has_trigger='.(array_key_exists('trigger', $data)?'yes':'no'));
    if (isset($data['trigger'])) {
        $t = $data['trigger'];
        $log('trigger='.(is_object($t)?get_class($t):gettype($t)));
        if (is_object($t)) $log('html='.$t->toHtml());
    }
    $attrs = $data['attributes'] ?? null;
    $log('attrs='.(is_object($attrs)?get_class($attrs):gettype($attrs)));
    if (is_object($attrs) && method_exists($attrs, 'all')) {
        $log('attrs_all_keys='.json_encode(array_keys($attrs->all())));
    }
} catch (Throwable $e) {
    $log('componentData EX: '.$e->getMessage());
}

echo file_get_contents($logFile);
echo "DONE\n";
