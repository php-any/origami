<?php
// fixtures/composer_read.phar is a native-PHP-generated SHA256-signed archive.
$result = require __DIR__ . '/fixtures/composer_read.phar';
if (strpos($result[0], 'phar://') !== 0
    || substr($result[0], -13) !== '/src/main.php'
    || $result[1] . '/main.php' !== $result[0]
    || $result[2] !== 'phar-fixture-data'
    || $result[3] !== true || $result[4] !== true) {
    throw new \Exception('PHAR read, stat or source identity failed: ' . json_encode($result));
}
if (file_get_contents('phar://origami-read.phar/src/data.txt') !== 'phar-fixture-data') {
    throw new \Exception('PHAR alias not available');
}
if (file_exists('phar://origami-read.phar/missing')) {
    throw new \Exception('Missing PHAR entry reported as existing');
}
echo "phar read passed\n";
