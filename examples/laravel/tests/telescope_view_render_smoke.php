<?php

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

// replicate helpers mix inline
function mix_inline($path, $manifestDirectory = '')
{
    $args = func_get_args();
    echo "func_get_args count: " . count($args) . "\n";
    echo "arg0 type: " . gettype($args[0]) . " val: " . $args[0] . "\n";
    $m = app(\Illuminate\Foundation\Mix::class);
    return $m($path, $manifestDirectory);
}

try {
    $r = mix_inline('app.css', 'vendor/telescope');
    echo "mix_inline ok: $r\n";
} catch (Throwable $e) {
    echo "mix_inline err: " . $e->getMessage() . "\n";
}

try {
    $r2 = mix('app.css', 'vendor/telescope');
    echo "mix helper ok: $r2\n";
} catch (Throwable $e) {
    echo "mix helper err: " . $e->getMessage() . "\n";
}
