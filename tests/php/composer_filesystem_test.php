<?php
$dir = sys_get_temp_dir() . '/origami-dir-' . uniqid();
mkdir($dir);
file_put_contents($dir . '/one.txt', 'contents');
try {
    $handle = opendir($dir);
    $entries = [];
    while (($name = readdir($handle)) !== false) { $entries[] = $name; }
    sort($entries);
    if ($entries !== ['.', '..', 'one.txt']) { throw new Exception('Directory entries mismatch'); }
    rewinddir($handle);
    if (readdir($handle) !== '.') { throw new Exception('Directory rewind failed'); }
    closedir($handle);
    $file = new SplFileInfo($dir . '/one.txt');
    if (!is_file($file) || is_dir($file) || is_link($file)) { throw new Exception('Stringable file status failed'); }
    $stat = lstat($file);
    if ($stat['size'] !== 8 || $stat[7] !== 8 || ($stat['mode'] & 0170000) !== 0100000) { throw new Exception('File stat mismatch'); }
    $stat = stat($dir);
    if (($stat['mode'] & 0170000) !== 0040000) { throw new Exception('Directory stat mismatch'); }
} finally {
    unlink($dir . '/one.txt');
    rmdir($dir);
}
echo "composer filesystem passed\n";
