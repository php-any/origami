<?php

$root = dirname(__DIR__, 2);
if (!is_dir($root . '/vendor')) {
    $root = getcwd();
}
require $root . '/vendor/autoload.php';
$app = require $root . '/bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

$compiler = app('blade.compiler');
$seen = [];

$compiler->directive('capture', function ($expression) use (&$seen) {
    $seen[] = $expression;
    echo "CAPTURE_EXPR=[" . var_export($expression, true) . "]\n";
    return RyanChandler\BladeCaptureDirective\BladeCaptureDirective::open($expression);
});

$path = $root . '/vendor/filament/filament/resources/views/components/logo.blade.php';
$contents = file_get_contents($path);
$compiled = $compiler->compileString($contents);
echo "SEEN_COUNT=" . count($seen) . "\n";
$pos = strpos($compiled, '(function ($args)');
echo "AROUND=[" . substr($compiled, max(0, $pos - 40), 80) . "]\n";
