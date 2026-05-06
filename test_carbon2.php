<?php
require __DIR__ . '/laravel/vendor/autoload.php';

// This is what Carbon::now() does:
$result = Carbon\Carbon::now();
echo "Carbon::now() = " . $result->format('Y-m-d H:i:s') . "\n";
