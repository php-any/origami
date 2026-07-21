<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Support\Collection;
use Illuminate\Support\Str;

$c = collect([1, 2, 3])->map(fn ($x) => $x * 2);
echo implode(',', $c->all()) . "\n";
echo Str::upper('hello') . "\n";
