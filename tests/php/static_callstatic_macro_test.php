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
}

echo MacroDemo::registerViaStatic(), PHP_EOL;
echo MacroDemo::hello(), PHP_EOL;
