<?php
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';

$r = new \Illuminate\Routing\Route(['GET'], '/', function () {});
try {
    $r->bind(\Illuminate\Http\Request::create('/'));
    $c = $r->getCompiled();
    echo "compiled=" . ($c === null ? 'null' : $c->getRegex()) . "\n";
} catch (\Throwable $e) {
    echo "error: " . $e->getMessage() . "\n";
}
