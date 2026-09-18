<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\View\AnonymousComponent;

$rm = new ReflectionMethod(AnonymousComponent::class, 'extractConstructorParameters');
$rm->setAccessible(true);
$params = $rm->invoke(null);
echo "params_type=".gettype($params)."\n";
echo "params_json=".json_encode($params)."\n";
echo "is_array=".(is_array($params)?'yes':'no')."\n";
echo "count=".count($params)."\n";
foreach ($params as $k => $v) {
    echo "p[$k]=".json_encode($v)." type=".gettype($v)."\n";
}

$flip = array_flip($params);
echo "flip_type=".gettype($flip)."\n";
echo "flip_json=".json_encode($flip)."\n";
echo "flip_count=".count($flip)."\n";
foreach ($flip as $k => $v) {
    echo "f[$k]=".json_encode($v)." type=".gettype($v)."\n";
}

$data = [
    'view' => 'filament-panels::components.page.simple',
    'data' => [],
];
echo "data_keys=".json_encode(array_keys($data))."\n";
$intersect = array_intersect_key($data, $flip);
echo "intersect_json=".json_encode($intersect)."\n";
echo "intersect_count=".count($intersect)."\n";
