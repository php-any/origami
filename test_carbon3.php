<?php
require __DIR__ . '/laravel/vendor/autoload.php';

$result = Illuminate\Support\Carbon::now();
echo "Illuminate\Support\Carbon::now() = " . $result->format('Y-m-d H:i:s') . "\n";
