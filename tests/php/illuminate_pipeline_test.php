<?php

/**
 * Illuminate\Pipeline\Pipeline 原生实现冒烟。未 AddClass 时跳过。
 */
if (!class_exists(\Illuminate\Pipeline\Pipeline::class, false)) {
    echo "skip: Pipeline 原生类未注册\n";
    return;
}

use Illuminate\Pipeline\Pipeline;

$passable = 'hello';
$result = (new Pipeline())
    ->send($passable)
    ->through([
        function ($p, $next) {
            return $next(strtoupper($p));
        },
        function ($p, $next) {
            return $next($p . '-world');
        },
    ])
    ->thenReturn();

if ($result !== 'HELLO-world') {
    echo "FAIL pipeline thenReturn expected HELLO-world got " . var_export($result, true) . "\n";
    exit(1);
}

echo "OK illuminate_pipeline_test\n";
