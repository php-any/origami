<?php
namespace tests\php;

trait Macroable {
    protected static $macros = [];
    public static function hasMacro($name) {
        return isset(static::$macros[$name]);
    }
    public static function macro($name, $macro) {
        static::$macros[$name] = $macro;
    }
    public function __call($method, $parameters) {
        if (! static::hasMacro($method)) {
            throw new \BadMethodCallException($method);
        }
        $macro = static::$macros[$method];
        return $macro(...$parameters);
    }
}

class BaseBag {
    use Macroable;
}

class ChildBag extends BaseBag {
    public static function hasMacro($name): bool {
        return parent::hasMacro($name) || BaseBag::hasMacro($name);
    }
}

echo "base: ";
var_dump(BaseBag::hasMacro('x'));
echo "child: ";
var_dump(ChildBag::hasMacro('x'));
BaseBag::macro('hello', fn() => 'hi');
echo "child after base macro: ";
var_dump(ChildBag::hasMacro('hello'));
echo "ok\n";
