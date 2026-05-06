<?php
require __DIR__ . '/laravel/vendor/autoload.php';

// Test calling setLastErrors on Carbon\Carbon
try {
    Carbon\Carbon::setLastErrors([]);
    echo "setLastErrors on Carbon\Carbon: OK\n";
} catch (Error $e) {
    echo "setLastErrors on Carbon\Carbon: " . $e->getMessage() . "\n";
}

// Test calling setLastErrors on Illuminate\Support\Carbon
try {
    Illuminate\Support\Carbon::setLastErrors([]);
    echo "setLastErrors on Illuminate\Support\Carbon: OK\n";
} catch (Error $e) {
    echo "setLastErrors on Illuminate\Support\Carbon: " . $e->getMessage() . "\n";
}
