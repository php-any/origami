<?php

require __DIR__.'/../../vendor/autoload.php';

$c = new Illuminate\Support\Collection(['a' => 1, 'b' => 2]);
echo 'class='.get_class($c)."\n";
echo 'exists_map='.(method_exists($c, 'map') ? 'yes' : 'no')."\n";
echo 'exists_mwk='.(method_exists($c, 'mapWithKeys') ? 'yes' : 'no')."\n";

try {
    $all = $c->all();
    echo 'all='.json_encode($all)."\n";
} catch (Throwable $e) {
    echo 'all_err='.$e->getMessage()."\n";
}

try {
    $mapped = $c->map(function ($v) { return $v; })->all();
    echo 'map='.json_encode($mapped)."\n";
} catch (Throwable $e) {
    echo 'map_err='.$e->getMessage()."\n";
}

$cb = function ($v, $k) { return [$k => $v * 10]; };
try {
    $got = $c->mapWithKeys($cb)->all();
    echo 'mwk='.json_encode($got)."\n";
} catch (Throwable $e) {
    echo 'mwk_err='.$e->getMessage()."\n";
}

$name = 'mapWithKeys';
try {
    $got2 = $c->$name($cb)->all();
    echo 'dyn='.json_encode($got2)."\n";
} catch (Throwable $e) {
    echo 'dyn_err='.$e->getMessage()."\n";
}
