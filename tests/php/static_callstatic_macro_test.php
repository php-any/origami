<?php

class MacroDemo
{
    protected static $macros = [];

    public static function __callStatic($method, $parameters)
    {
        if ($method === 'macro') {
            static::$macros[$parameters[0]] = $parameters[1];
            return 'registered';
        }
        if (isset(static::$macros[$method])) {
            return (static::$macros[$method])(...$parameters);
        }
        throw new BadMethodCallException($method);
    }

    protected static function registerViaStatic()
    {
        return static::macro('hello', fn () => 'world');
    }

    public static function register() { return static::registerViaStatic(); }
}

try {
    MacroDemo::registerViaStatic();
    throw new RuntimeException('inaccessible method must reach __callStatic');
} catch (BadMethodCallException $e) {
    if ($e->getMessage() !== 'registerViaStatic') { throw $e; }
}
echo MacroDemo::register(), PHP_EOL;
echo MacroDemo::hello(), PHP_EOL;
