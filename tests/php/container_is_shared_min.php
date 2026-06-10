<?php

$c = \Container\Container::getInstance();
$abstract = 'App\\Http\\Kernel';

if ($c->isShared($abstract)) {
    echo "isShared=1\n";
} else {
    echo "isShared=0\n";
}
