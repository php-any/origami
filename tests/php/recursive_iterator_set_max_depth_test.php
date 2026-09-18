<?php

$it = new RecursiveIteratorIterator(new RecursiveDirectoryIterator(__DIR__), RecursiveIteratorIterator::SELF_FIRST);
$it->setMaxDepth(0);
echo 'max='.$it->getMaxDepth(), PHP_EOL;
$n = 0;
foreach ($it as $f) {
    $n++;
    if ($n > 5) break;
}
echo 'ok count>='.$n, PHP_EOL;
