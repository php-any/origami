<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "A\n";
$v1 = view('origami-capture-test', ['x' => 1]);
echo "B class=".get_class($v1)."\n";
echo "C isset=".(isset($v1->layoutConfig)?'yes':'no')."\n";
echo "D\n";
$r = $v1->layoutConfig ?? 'DEFAULT';
echo "E r=".$r."\n";

// Reflect __isset / data
$ref = new ReflectionClass($v1);
$prop = $ref->getProperty('data');
$prop->setAccessible(true);
$data = $prop->getValue($v1);
echo "data_keys=".implode(',', array_keys($data))."\n";
echo "data_has_lc=".(array_key_exists('layoutConfig', $data)?'yes':'no')."\n";

echo "DONE\n";
