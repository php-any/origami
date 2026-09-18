<?php
trait Macroable {
    protected static $macros = [];
    public static function hasMacro($name) {
        return isset(static::$macros[$name]);
    }
    public static function macro($name, $macro) {
        static::$macros[$name] = $macro;
    }
    public static function __callStatic($method, $parameters) {
        if (! static::hasMacro($method)) {
            throw new BadMethodCallException($method);
        }
        $macro = static::$macros[$method];
        return $macro(...$parameters);
    }
    public function __call($method, $parameters) {
        if (! static::hasMacro($method)) {
            throw new BadMethodCallException($method);
        }
        $macro = static::$macros[$method];
        return $macro(...$parameters);
    }
}
class Foo {
    use Macroable;
}
echo "hasMacro: ";
var_dump(Foo::hasMacro('x'));
Foo::macro('hello', fn() => 'world');
echo Foo::hello(), PHP_EOL;
$f = new Foo();
try {
  $f->missing();
} catch (Throwable $e) {
  echo "missing: ".$e->getMessage(), PHP_EOL;
}
