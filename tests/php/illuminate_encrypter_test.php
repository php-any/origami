<?php

if (!class_exists(\Illuminate\Encryption\Encrypter::class, false)) {
    echo "skip: Encrypter 原生类未注册\n";
    return;
}

$key = \Illuminate\Encryption\Encrypter::generateKey('aes-256-cbc');
if (!is_string($key) || strlen($key) !== 32) {
    echo "FAIL generateKey len\n";
    exit(1);
}

$e = new \Illuminate\Encryption\Encrypter($key, 'aes-256-cbc');
if ($e->getKey() !== $key) {
    echo "FAIL getKey\n";
    exit(1);
}

$enc = $e->encryptString('hello');
if (!\Illuminate\Encryption\Encrypter::appearsEncrypted($enc)) {
    echo "FAIL appearsEncrypted\n";
    exit(1);
}
if ($e->decryptString($enc) !== 'hello') {
    echo "FAIL decryptString\n";
    exit(1);
}

$prev = \Illuminate\Encryption\Encrypter::generateKey('aes-256-cbc');
$e->previousKeys([$prev]);
$all = $e->getAllKeys();
if (!is_array($all) || count($all) < 2) {
    echo "FAIL getAllKeys\n";
    exit(1);
}

echo "OK illuminate_encrypter_test\n";
