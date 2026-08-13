<?php

require_once __DIR__ . '/../../examples/laravel/vendor/autoload.php';

use Illuminate\Support\Arr;

$data = ['view' => ['paths' => ['/tmp/views']]];
$v = Arr::get($data, 'view.paths.0');
if ($v !== '/tmp/views') {
    echo "FAIL Arr::get dot numeric: ";
    var_export($v);
    echo "\n";
    exit(1);
}

echo "PASS\n";
