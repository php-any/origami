<?php

function in_identity($x) { return $x; }

class InWrap {
    public function __call($m, $a) { return null; }
}
function in_opt($x) { return new InWrap(); }

// chain after identity(null) - Wrap always returned from opt
$c = ['v' => in_opt(null)->m()];
var_dump($c);

// chain directly on identity(null) return value - null->m() would fail at runtime
try {
    $d = ['v' => in_identity(null)->m()];
    var_dump($d);
} catch (\Throwable $e) {
    Log::info('identity null chain error expected');
}

// parse test: optional-like with null inside parens then chain on wrapper
$e = ['v' => in_opt(in_identity(null))->m()];
var_dump($e);

Log::info('chain tests done');
