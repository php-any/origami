<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Hashing\BcryptHasher;

$hasher = new BcryptHasher(['rounds' => 4]);
$hash = $hasher->make('secret');

if (!is_string($hash) || $hash === '') {
    echo "FAIL: make returned empty\n";
    exit(1);
}

if (!$hasher->check('secret', $hash)) {
    echo "FAIL: check should pass\n";
    exit(1);
}

if ($hasher->check('wrong', $hash)) {
    echo "FAIL: check should fail for wrong password\n";
    exit(1);
}

echo "PASS\n";
