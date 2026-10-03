<?php
function ser_check($ok, $message) { if (!$ok) throw new Exception($message); }
class SerBase { private $hidden = 'base'; protected $shared = null; public $base = 1; }
class SerChild extends SerBase { private $hidden = 'child'; public $nil = null; public int $uninitialized; public static $static = 8; public function values() { return [$this->hidden, $this->shared, $this->nil]; } }
$object = new SerChild();
$encoded = serialize($object);
echo bin2hex($encoded), "\n";
$copy = unserialize($encoded);
ser_check($copy instanceof SerChild && $copy !== $object && $copy->values() === ['child', null, null], 'visibility/default values');
ser_check(serialize($copy) === $encoded, 'object round trip');
$ref = 7;
$graph = [$object, $object, &$ref, &$ref, ['x' => 1], ['x' => 1]];
$encoded = serialize($graph);
echo bin2hex($encoded), "\n";
$copy = unserialize($encoded);
ser_check($copy[0] === $copy[1] && $copy[0] !== $object, 'object identity');
$copy[2] = 9;
ser_check($copy[3] === 9 && $ref === 7, 'scalar references');
$copy[4]['x'] = 2;
ser_check($copy[5]['x'] === 1, 'array value copies');
$cycle = new stdClass(); $cycle->self = $cycle;
$decoded = unserialize(serialize($cycle));
ser_check($decoded === $decoded->self, 'object cycle');
$recursive = []; $recursive['self'] = &$recursive;
echo bin2hex(serialize($recursive)), "\n";
$recursiveCopy = unserialize(serialize($recursive));
$recursiveCopy['self']['changed'] = 3;
ser_check($recursiveCopy['self']['self']['changed'] === 3, 'recursive array references');
class SerHooks {
    private int $value = 5;
    public function __serialize(): array { return [0 => $this->value, 'nil' => null]; }
    public function __unserialize(array $data): void { $this->value = $data[0] + 1; }
    public function value(): int { return $this->value; }
    public function __wakeup(): void { throw new Exception('wakeup ran with unserialize'); }
}
$hooks = unserialize(serialize(new SerHooks()));
ser_check($hooks->value() === 6, 'serialization hooks');
class SerWake { public $value = 1; public function __wakeup(): void { $this->value++; } }
ser_check(unserialize(serialize(new SerWake()))->value === 2, 'wakeup');
class SerTyped { public int $value; }
$typed = unserialize('O:8:"SerTyped":1:{s:5:"value";i:4;}');
ser_check($typed->value === 4, 'typed initialization');
try { unserialize('O:8:"SerTyped":1:{s:5:"value";s:1:"4";}'); throw new Exception('typed string accepted'); } catch (TypeError $e) {}
$incomplete = unserialize($encoded, ['allowed_classes' => false]);
ser_check(get_class($incomplete[0]) === '__PHP_Incomplete_Class', 'allowed_classes');
ser_check(serialize($incomplete) === $encoded, 'incomplete class round trip');
$floats = [1.25, 1000000.0, 1.0e17, 1.0e-5];
echo serialize($floats), "\n";
ser_check(unserialize(serialize($floats)) === $floats, 'float round trip');
try { serialize(function () {}); throw new Exception('serialized closure'); } catch (Exception $e) { ser_check($e->getMessage() === "Serialization of 'Closure' is not allowed", 'closure rejection'); }
set_error_handler(function () { return true; });
ser_check(unserialize('a:999999999:{') === false, 'malformed size');
ser_check(unserialize('s:99:"x";') === false, 'invalid byte length');
restore_error_handler();
echo "serialization graph contract OK\n";
