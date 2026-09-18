<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$bag = new Illuminate\View\ComponentAttributeBag(['class' => 'probe']);

echo "has_method_class=".(method_exists($bag, 'class')?'yes':'no')."\n";
echo "methods_has=".(in_array('class', get_class_methods($bag))?'yes':'no')."\n";

$wrapped = Illuminate\Support\Arr::wrap(['x']);
echo "wrap_type=".gettype($wrapped)."\n";
echo "wrap=".json_encode($wrapped)."\n";

$css = Illuminate\Support\Arr::toCssClasses(['x']);
echo "css=".var_export($css, true)."\n";

$css2 = Illuminate\Support\Arr::toCssClasses(['fi-logo' => true, 'x' => false]);
echo "css2=".var_export($css2, true)."\n";

$merged = $bag->merge(['class' => 'x']);
echo "merge_direct=".$merged."\n";

// Call class via reflection to bypass dispatch
$ref = new ReflectionMethod($bag, 'class');
$out = $ref->invoke($bag, ['x']);
echo "via_reflection=".$out."\n";

echo "via_call=".$bag->class(['x'])."\n";
