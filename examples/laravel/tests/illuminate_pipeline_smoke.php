<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Pipeline\Pipeline;

$result = (new Pipeline())
    ->send(1)
    ->through([
        function ($passable, $next) {
            return $next($passable + 1);
        },
        function ($passable, $next) {
            return $next($passable * 2);
        },
    ])
    ->then(fn ($passable) => $passable + 10);

if ($result !== 14) {
    echo "FAIL: expected 14, got {$result}\n";
    exit(1);
}

echo "PASS\n";
