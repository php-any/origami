<?php
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';

$r = new \Illuminate\Routing\Route(['GET'], '/', function () {});
$request = \Illuminate\Http\Request::create('/');
try {
    $ok = $r->matches($request);
    $c = $r->getCompiled();
    echo "matches=" . ($ok ? '1' : '0') . " compiled=" . ($c === null ? 'null' : 'ok') . "\n";
} catch (\Throwable $e) {
    echo "error: " . $e->getMessage() . "\n";
}
