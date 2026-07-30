<?php

class Foo
{
    public static function __callStatic(string $method, array $parameters): mixed
    {
        echo "called=[$method] args=" . json_encode($parameters) . "\n";
        return "ok:$method";
    }

    public static function bar()
    {
        return 'bar';
    }
}

$method = 'bar';
echo Foo::$method() . "\n";

$method = 'baz';
echo Foo::$method('x', 1) . "\n";

$class = 'Foo';
$method = 'bar';
echo $class::$method() . "\n";
