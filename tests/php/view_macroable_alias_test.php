<?php
namespace tests\php;

trait Macroable {
    protected static $macros = [];
    public static function hasMacro($name) {
        echo "hasMacro body\n";
        return isset(static::$macros[$name]);
    }
    public static function __callStatic($method, $parameters) {
        echo "callStatic $method\n";
        if (! static::hasMacro($method)) {
            throw new \BadMethodCallException($method);
        }
        return null;
    }
    public function __call($method, $parameters) {
        echo "call $method\n";
        if (! static::hasMacro($method)) {
            throw new \BadMethodCallException($method);
        }
        return $this->macroCall($method, $parameters);
    }
}

class View {
    use Macroable {
        __call as macroCall;
    }
    public function __call($method, $parameters) {
        echo "View::__call $method\n";
        if (static::hasMacro($method)) {
            return $this->macroCall($method, $parameters);
        }
        throw new \BadMethodCallException($method);
    }
}

echo "direct: ";
var_dump(View::hasMacro('x'));
$v = new View();
try { $v->foo(); } catch (\Throwable $e) { echo "ex: ".$e->getMessage()."\n"; }
