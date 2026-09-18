<?php
/**
 * Arr::map / array_values 返回类型。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$items = ['App\\A', 'App\\B'];
echo "array_values type=".gettype(array_values($items))." json=".json_encode(array_values($items))."\n";

$mapped = Illuminate\Support\Arr::map($items, fn ($w, $k) => $w);
echo "Arr::map type=".gettype($mapped)." countable=".(is_countable($mapped)?'y':'n')." json=".json_encode($mapped)."\n";
echo "is_array=". (is_array($mapped)?'y':'n')."\n";

$c = collect($items)->values()->map(fn ($w, $k) => $w);
echo "collection class=".get_class($c)."\n";
$all = $c->all();
echo "all type=".gettype($all)." is_array=".(is_array($all)?'y':'n')." json=".json_encode($all)."\n";

// 直接看 items
$ref = new ReflectionClass($c);
$p = $ref->getProperty('items');
$p->setAccessible(true);
$raw = $p->getValue($c);
echo "items type=".gettype($raw)." is_array=".(is_array($raw)?'y':'n')." json=".json_encode($raw)."\n";

echo "DONE\n";
