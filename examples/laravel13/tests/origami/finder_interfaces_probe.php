<?php

require __DIR__.'/../../vendor/autoload.php';

use Symfony\Component\Finder\Finder;

$f = Finder::create();
$rc = new ReflectionClass($f);
echo 'class: '.$rc->getName().PHP_EOL;
echo 'interfaces:'.PHP_EOL;
foreach ($rc->getInterfaces() as $iface) {
    echo ' - '.$iface->getName().PHP_EOL;
}
echo 'implements IteratorAggregate directly: '.($rc->implementsInterface(IteratorAggregate::class) ? 'yes' : 'no').PHP_EOL;
