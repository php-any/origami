<?php
$path = __DIR__ . '/../../vendor/filament/support/src/SupportServiceProvider.php';
$src = file_get_contents($path);
// Extract the exact return string from source between return " and ";
$start = strpos($src, "Blade::directive('capture'");
$chunk = substr($src, $start, 800);
echo "CHUNK:\n", $chunk, "\n---\n";

// Hex dump of the return string portion
$rpos = strpos($chunk, 'return "');
$rchunk = substr($chunk, $rpos, 400);
echo "HEX:\n";
for ($i = 0; $i < strlen($rchunk) && $i < 200; $i++) {
    $c = $rchunk[$i];
    $o = ord($c);
    if ($o < 32 || $o > 126) {
        echo sprintf('\\x%02x', $o);
    } else {
        echo $c;
    }
}
echo "\n";
