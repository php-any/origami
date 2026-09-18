<?php

class Foo
{
    public function bar(): void {}
}

$m = (new ReflectionClass(Foo::class))->getMethod('bar');
echo $m->name, PHP_EOL;
echo $m->class, PHP_EOL;

$fn = function (): self {
    return $this;
};
echo 'closure self ok', PHP_EOL;
