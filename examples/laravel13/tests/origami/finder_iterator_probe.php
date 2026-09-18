<?php

require __DIR__.'/../../vendor/autoload.php';

use Symfony\Component\Finder\Finder;

$f = Finder::create()->files()->in(__DIR__)->depth(0)->sortByName();
$it = $f->getIterator();
echo 'iterator class: '.get_class($it).PHP_EOL;
echo 'is Iterator: '.(is_a($it, Iterator::class) ? 'yes' : 'no').PHP_EOL;
echo 'is Traversable: '.(is_a($it, Traversable::class) ? 'yes' : 'no').PHP_EOL;
echo 'has valid: '.(method_exists($it, 'valid') ? 'yes' : 'no').PHP_EOL;

try {
    $arr = iterator_to_array($it, false);
    echo 'direct iterator count: '.count($arr).PHP_EOL;
} catch (Throwable $e) {
    echo 'direct iterator error: '.$e->getMessage().PHP_EOL;
}

try {
    $arr = iterator_to_array($f, false);
    echo 'finder count: '.count($arr).PHP_EOL;
} catch (Throwable $e) {
    echo 'finder error: '.$e->getMessage().PHP_EOL;
}
