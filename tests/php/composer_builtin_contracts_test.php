<?php
if (extension_loaded('xdebug') || extension_loaded('origami_nonexistent_extension')) {
    throw new \Exception('Unimplemented extension reported as loaded');
}
if (!extension_loaded('Phar') || !class_exists('Phar')) {
    throw new \Exception('Phar capability missing');
}
if (!defined('PHP_BINARY') || !is_file(PHP_BINARY)) {
    throw new \Exception('PHP_BINARY must identify the running executable');
}
if (DIRECTORY_SEPARATOR === '\\' && !defined('PHP_WINDOWS_VERSION_BUILD')) {
    throw new \Exception('Windows platform constants missing');
}
$file = tempnam(sys_get_temp_dir(), 'origami-composer-');
try {
    $write = 'file_put_contents';
    if ($write(...[$file, 'test']) !== 4 || file_get_contents($file) !== 'test') {
        throw new \Exception('Optional stream context became required');
    }
} finally {
    unlink($file);
}
echo "composer builtin contracts passed\n";
