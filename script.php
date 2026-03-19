<?php

class A {}

$var[0] = new A();


$class = is_object($var[0]) ? get_class($var[0]) : $var[0];

var_dump($class);

if (! class_exists($class)) {
    return false;
}