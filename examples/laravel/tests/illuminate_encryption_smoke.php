<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Encryption\Encrypter;

$key = str_repeat('a', 16); // aes-128-cbc
$encrypter = new Encrypter($key, 'aes-128-cbc');

$cipher = $encrypter->encryptString('hello');
$plain = $encrypter->decryptString($cipher);

if ($plain !== 'hello') {
    echo "FAIL: expected hello, got {$plain}\n";
    exit(1);
}

echo "PASS\n";
