<?php

use Illuminate\Support\Collection;

require __DIR__.'/../../vendor/autoload.php';

$all = (new Collection([]))->map(fn ($v) => 1)->all();
echo "gettype=".gettype($all)."\n";
echo "class=".(is_object($all) ? get_class($all) : 'n/a')."\n";
echo "is_array=".var_export(is_array($all), true)."\n";
echo "bool=".var_export((bool)$all, true)."\n";
echo "var_export=".var_export($all, true)."\n";

$c = new Collection([]);
$ref = new ReflectionClass($c);
// skip reflection if broken - dump items via toArray
echo "toArray_type=".gettype($c->toArray())." toArray_bool=".var_export((bool)$c->toArray(), true)."\n";

$plain = [];
echo "plain_type=".gettype($plain)." plain_bool=".var_export((bool)$plain, true)."\n";
