<?php
namespace tests\php;

trait Macroable {
    protected static $macros = [];
    public static function hasMacro($name) {
        return isset(static::$macros[$name]);
    }
    public static function __callStatic($method, $parameters) {
        echo "__callStatic:$method\n";
        if (! static::hasMacro($method)) {
            throw new \BadMethodCallException($method);
        }
        return null;
    }
    public function __call($method, $parameters) {
        echo "__call:$method\n";
        if (! static::hasMacro($method)) {
            throw new \BadMethodCallException($method);
        }
        return null;
    }
}

class View {
    use Macroable {
        __call as macroCall;
    }
    public function __call($method, $parameters) {
        echo "View::__call:$method\n";
        if (static::hasMacro($method)) {
            return $this->macroCall($method, $parameters);
        }
        throw new \BadMethodCallException($method);
    }
}

$v = new View();
try { $v->missing(); } catch (\Throwable $e) { echo "ex: ".$e->getMessage()."\n"; }

// Force many static::hasMacro from instance context
class Probe {
    use Macroable;
    public function check() {
        for ($i = 0; $i < 5; $i++) {
            var_dump(static::hasMacro('x'));
        }
    }
}
(new Probe())->check();
