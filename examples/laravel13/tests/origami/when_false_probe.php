<?php
require __DIR__.'/../../vendor/autoload.php';
foreach ([false, 0, null, ''] as $v) {
    try {
        $r = collect([1])->when($v, fn ($c) => $c->map(fn ($x) => $x * 10));
        echo gettype($v).' ok count='.$r->count()."\n";
    } catch (Throwable $e) {
        echo gettype($v).' err='.$e->getMessage()."\n";
    }
}
