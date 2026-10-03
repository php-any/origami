<?php

class Foo
{
    public function bar(): void {}
    public function closure(): Closure {
        return function (): self { return $this; };
    }
}

$m = (new ReflectionClass(Foo::class))->getMethod('bar');
echo $m->name, PHP_EOL;
echo $m->class, PHP_EOL;

$foo = new Foo();
$fn = $foo->closure();
if ($fn() !== $foo) { throw new RuntimeException('closure self identity'); }
echo 'closure self ok', PHP_EOL;
